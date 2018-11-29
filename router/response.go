package router

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ContentType int

var (
	ContentTypeNon  ContentType = 0
	ContentTypeJson ContentType = 1
	ContentTypeXml  ContentType = 2
	ContentTypeByte ContentType = 3
)

func (data *CommRouter) DealWithRequest() {
	data.ContentType = defaultContentType
	if values, ok := data.Req.Header["Content-Type"]; ok {
		for _, ct := range values {
			t := strings.ToLower(ct)
			if t == "application/json" {
				data.ContentType = ContentTypeJson
			} else if t == "application/xml" || t == "text/xml" {
				data.ContentType = ContentTypeXml
			}
		}
	} else if defaultContentType == ContentTypeNon {
		data.Err = errors.New("Header don't have Content-Type")
		data.Message = data.Err.Error()
		data.ErrCode = http.StatusBadRequest
		return
	}

	if data.Req.Method == "POST" {
		// read request body to byte
		if body, err := ioutil.ReadAll(data.Req.Body); err != nil || len(body) == 0 {
			data.Err = err
			data.ErrCode = http.StatusBadRequest
			if err == nil {
				data.Err = errors.New("body is null")
			}
			data.Message = "Read Request Body Fail: " + data.Err.Error()
			return
		} else {
			data.ReqBody = body

			if data.ContentType == ContentTypeJson {
				json.Unmarshal(body, &data.ReqMap)
			} else if data.ContentType == ContentTypeXml {
				xml.Unmarshal(body, (*XmlMap)(&data.ReqMap))
			}
		}
	}
}

const (
	XmlErrResponseStr  = `<xml><code>%v</code><message>%v</message><uid>%s</uid></xml>`
	JsonErrResponseStr = `{"code":"%v", "message":"%v", "uid":"%s"}`
)

func (data *CommRouter) DealWithResponse() {
	// check function error
	if data.Err != nil {
		if len(data.RspBody) == 0 {
			if data.ErrCode == nil {
				data.ErrCode = "FAIL"
			}
			if data.Message == "" {
				data.Message = data.Err.Error()
			}
			if data.ContentType == ContentTypeJson {
				data.RspBody = []byte(fmt.Sprintf(JsonErrResponseStr, data.ErrCode, data.Message, data.Uid))
			} else if data.ContentType == ContentTypeXml {
				data.RspBody = []byte(fmt.Sprintf(XmlErrResponseStr, data.ErrCode, data.Message, data.Uid))
			} else {
				data.RspBody = []byte(data.Message)
			}
		}
		return
	} else if len(data.RspBody) > 0 {
		data.CheckStatus(http.StatusOK)
		return
	}

	switch data.ContentType {
	case ContentTypeJson:
		if data.Result != nil {
			// convert data.Result struct to byte
			if data.RspBody, data.Err = json.Marshal(data.Result); data.Err == nil {
				data.CheckStatus(http.StatusOK)
			} else {
				data.RspBody = []byte(fmt.Sprintf(JsonErrResponseStr, "FAIL", data.Err, data.Uid))
			}
		} else if len(data.RspMap) > 0 {
			// convert data.Result struct to byte
			if data.RspBody, data.Err = json.Marshal(data.RspMap); data.Err == nil {
				data.CheckStatus(http.StatusOK)
			} else {
				data.RspBody = []byte(fmt.Sprintf(JsonErrResponseStr, "FAIL", data.Err, data.Uid))
			}
		}
	case ContentTypeXml:
		if data.Result != nil {
			// convert data.Result struct to byte
			if data.RspBody, data.Err = xml.Marshal(data.Result); data.Err == nil {
				data.CheckStatus(http.StatusOK)
			} else {
				data.RspBody = []byte(fmt.Sprintf(XmlErrResponseStr, "FAIL", data.Err, data.Uid))
			}
		} else if len(data.RspMap) > 0 {
			// convert data.RspMap
			if data.RspBody, data.Err = xml.Marshal(XmlMap(data.RspMap)); data.Err == nil {
				data.CheckStatus(http.StatusOK)
			} else {
				data.RspBody = []byte(fmt.Sprintf(XmlErrResponseStr, "FAIL", data.Err, data.Uid))
			}
		}
	}

	return
}

func (data *CommRouter) CheckStatus(status int) {
	if data.Status == 0 {
		data.Status = status
	}
}

// Convert xml to map[string]string and map[string]string to xml function
type XmlMap map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (m XmlMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}
	start.Name.Local = "xml"
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func (m *XmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}
