package router

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

// request count
var count int32

// Server common struct
type CommRouter struct {
	// use ext values
	Ext map[string]interface{}

	// response data
	Err     error             // error
	ErrCode interface{}       // can use int or string
	Message string            // err not null, response message
	RspBody []byte            // response body will read RspBody -> Result -> RspMap
	Result  interface{}       // json or xml struct data to response
	RspMap  map[string]string // json or xml map to response

	// requset data
	req     *http.Request
	reqMap  map[string]string // request body convert to map
	rmc     bool              // check request map
	body    []byte            // request body ([]byte)
	bdc     bool              // check read body
	uid     string            // request uid
	start   time.Time         // request deal with start time
	content ContentType       // request content type
	ctc     bool              // check content-type

	// response write data
	header  map[string]string // response header add set
	status  int               // http.Status
	resConv bool              // make response convert json or xml
}

func (c *CommRouter) GetRequest() *http.Request {
	return c.req
}

// read request body where method is post
// check body bytes lenght less than max setting
func (c *CommRouter) GetBody() ([]byte, error) {
	if c.req.Method != "POST" {
		return nil, errors.New("Request method not POST")
	}

	if c.bdc {
		return c.body, nil
	} else {
		c.bdc = true
		if c.req.ContentLength > getMaxBodyRead() {
			return nil, errors.New("Body bytes more than max setting")
		} else {
			c.body = make([]byte, c.req.ContentLength)
			c.req.Body.Read(c.body)
			return c.body, nil
		}
	}
}

func (c *CommRouter) GetUid() string {
	return c.uid
}

func (c *CommRouter) GetRequestTime() time.Time {
	return c.start
}

func (c *CommRouter) GetContentType() ContentType {
	if c.ctc {
		return c.content
	} else {
		c.ctc = true
		c.content = ContentTypeNon
		if values, ok := c.req.Header["Content-Type"]; ok {
			for _, ct := range values {
				t := strings.ToLower(ct)
				if strings.Contains(t, "json") {
					c.content = ContentTypeJson
				} else if strings.Contains(t, "xml") {
					c.content = ContentTypeXml
				} else if strings.Contains(t, "text/plain") {
					c.content = ContentTypeText
				} else if strings.Contains(t, "multipart/form-data") {
					c.content = ContentTypeData
				} else if strings.Contains(t, "text/html") {
					c.content = ContentTypeHtml
				}
			}
		}
		return c.content
	}
}

func (c *CommRouter) GetRequestMap() (map[string]string, error) {
	if c.rmc {
		return c.reqMap, nil
	} else {
		c.rmc = true
		body, err := c.GetBody()
		if err != nil {
			return nil, err
		}
		content := c.GetContentType()

		if content == ContentTypeJson {
			var result map[string]interface{}
			if err = json.Unmarshal(body, &result); err == nil {
				for k, v := range result {
					c.reqMap[k] = fmt.Sprint(v)
				}
			} else {
				return nil, err
			}
		} else if content == ContentTypeXml {
			err = xml.Unmarshal(body, (*XmlMap)(&c.reqMap))
		} else {
			return nil, ErrContentTypeUnknown
		}
		return c.reqMap, err
	}
}

func (c *CommRouter) PutError(args ...interface{}) {
	c.Err = errors.New(argsToStr(args))
}

func (c *CommRouter) PutMessage(message ...interface{}) {
	c.Message = argsToStr(message)
}

// convert response map or interface to json or xml
// according to request content-type [ContentTypeJson | ContentTypeXml]
// when server set default convert type, will use than
// default false, not open
func (c *CommRouter) SetResponseConvert() {
	c.resConv = true
}

func (c *CommRouter) SetHeader(key, value string) {
	c.header[key] = value
}

func (c *CommRouter) SetStatus(status int) {
	c.status = status
}

func argsToStr(args []interface{}) string {
	var result = make([]string, len(args))
	for i, v := range args {
		result[i] = fmt.Sprint(v)
	}
	return strings.Join(result, " ")
}

// Common Deal With Request and Response
func CommonDealWith(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// requset count
	atomic.AddInt32(&count, 1)
	defer atomic.AddInt32(&count, -1)

	var timestamp = time.Now()
	comm := &CommRouter{
		req:    r,
		reqMap: make(map[string]string),
		start:  timestamp,
		uid:    MakeUid(r, timestamp),
		Ext:    make(map[string]interface{}),
		RspMap: make(map[string]string),
		header: make(map[string]string),
	}

	// make request log
	log.Printf("\t%s\t%s\t%s", r.RemoteAddr, r.URL.Path, comm.uid)

	// Deal with functions
	if fs := getFunction(r.URL.Path); fs == nil {
		comm.Err = errors.New("Url Path Error")
		comm.Message = comm.Err.Error()
	} else {
		if fs.method != r.Method {
			comm.PutError("Request method not right:", r.Method)
			comm.PutMessage("Request method should be:", fs.method)
		} else {
			for _, f := range fs.function {
				if f(comm); comm.Err != nil {
					break
				}
			}
		}
	}

	// convert response data
	if comm.resConv {
		// Deal with response
		comm.dealWithResponse()
	}

	// response header write
	if len(comm.header) > 0 {
		for k, v := range comm.header {
			w.Header().Add(k, v)
		}
	}
	// write response status
	if comm.status != 0 {
		w.WriteHeader(comm.status)
	}
	// return response body
	if len(comm.RspBody) > 0 {
		w.Write(comm.RspBody)
	}

	// write end log
	go requestAndResponseLog(comm)

	return
}

// make request and response log
// latency is event from request to response need time, (ms)
func requestAndResponseLog(comm *CommRouter) {
	var latency = time.Now().Sub(comm.start).Seconds() * 1000
	var data = map[string]interface{}{
		"uid":      comm.uid,
		"name":     "request_and_response",
		"date":     comm.start.Format("2006-01-02 15:04:05"),
		"status":   comm.status,
		"latency":  latency,
		"ip":       comm.req.RemoteAddr,
		"method":   comm.req.Method,
		"path":     comm.req.URL.Path,
		"request":  bytesToString(comm.body),
		"response": bytesToString(comm.RspBody),
	}
	log.Println(data)
}

func bytesToString(data []byte) string {
	temp := strings.Replace(string(data), "\n", "", -1)
	return strings.Replace(temp, "\t", " ", -1)
}
