package router

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//
type CommRouter struct {
	// requset data
	Req         *http.Request
	ReqMap      map[string]string
	ReqBody     []byte
	JobId       string
	StartTime   time.Time
	ContentType ContentType

	// response data
	Err     error
	ErrCode interface{} // can use int or string
	Message string
	Status  int
	RspBody []byte
	Result  interface{}
	RspMap  map[string]string
	EndTime time.Time
}

// Common Deal With Request and Response
func CommonDealRequest(w http.ResponseWriter, r *http.Request) {
	var timestamp = time.Now()
	comm := &CommRouter{
		Req:       r,
		ReqMap:    make(map[string]string),
		StartTime: timestamp,
		JobId:     MakeUid(r, timestamp),
	}

	log.Printf("\t%s\t%s\t%s", r.RemoteAddr, r.URL.Path, comm.JobId)

	// Deal with request
	comm.DealWithRequest()

	if comm.Err == nil {

		// result is interface
		// when err not equal null, return fail and message
		if f := GetFunction(r.URL.Path); f == nil {
			comm.Err = errors.New("Url Path Error")
			comm.Message = comm.Err.Error()
		} else {
			f.function(comm)
		}
	}

	// Deal with response
	comm.DealWithResponse()

	w.WriteHeader(comm.Status)
	w.Write(comm.RspBody)

	// write end log
	comm.EndTime = time.Now()
	go requestAndResponseLog(comm)

	return
}

func requestAndResponseLog(comm *CommRouter) {
	var latency = time.Now().Sub(comm.StartTime).Seconds()
	var data = map[string]interface{}{
		"uid":     comm.JobId,
		"name":    "request_and_response",
		"date":    comm.StartTime.Format("2006/01/02 15:04:05"),
		"status":  comm.Status,
		"latency": latency,
		"ip":      comm.Req.RemoteAddr,
		"method":  comm.Req.Method,
		"path":    comm.Req.URL.Path,
	}
	if comm.Req.Method == "POST" && len(comm.ReqBody) < 2048 {
		req := strings.Replace(string(comm.ReqBody), "\n", "", -1)
		data["request"] = strings.Replace(req, "\t", "", -1)
	}
	if len(comm.RspBody) > 0 && len(comm.RspBody) < 2048 {
		rsp := strings.Replace(string(comm.RspBody), "\n", "", -1)
		data["response"] = strings.Replace(rsp, "\t", "", -1)
	}
	log.Println(data)
}

func ByteToString(data []string) string {
	var result = make([]byte, 0)
	for _, row := range data {
		d, _ := strconv.Atoi(row)
		result = append(result, byte(d/26+48), byte(d%26+65))
	}
	return string(result)
}

func MakeUid(r *http.Request, timestamp time.Time) string {
	ip := strings.Split(strings.Split(r.RemoteAddr, ":")[0], ".")

	rand.Seed(timestamp.UnixNano())
	var bte = make([]byte, 0, 10)
	for i := 0; i < 10; i++ {
		bte = append(bte, byte(rand.Intn(26)+65))
	}
	// UID = IP + Time + Range
	return fmt.Sprintf("%s-%s-%s",
		ByteToString(ip),
		timestamp.Format("060102150405"),
		string(bte),
	)
}
