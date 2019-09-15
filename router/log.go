package router

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type RequestAndResponse struct {
	Time     string
	Name     string
	Uid      string
	Ip       string
	Path     string
	Method   string
	Date     string
	Latency  string
	Status   int
	Request  string
	Response string
}

// make request and response log
// latency is event from request to response need time, (ms)
func requestAndResponseLog(comm *CommRouter) {
	var data = &RequestAndResponse{
		Time:    comm.start.Format("2006-01-02 15:04:05"),
		Name:    "RequestAndResponse",
		Uid:     comm.uid,
		Ip:      comm.req.RemoteAddr,
		Path:    comm.req.URL.Path,
		Method:  comm.req.Method,
		Date:    comm.start.Format("2006-01-02"),
		Status:  comm.status,
		Latency: fmt.Sprintf("%dms", time.Now().Sub(comm.start).Nanoseconds()/1000000),
	}
	if comm.req.Method == "POST" {
		data.Request = bytesToString(comm.body)
	} else if comm.req.Method == "GET" && comm.req.URL.ForceQuery {
		data.Request = comm.req.URL.RawQuery
	}
	if len(comm.RspBody) > 0 {
		data.Response = bytesToString(comm.RspBody)
	}
	// TODO: log write in file

	log.Println(data)
}

func bytesToString(data []byte) string {
	temp := strings.Replace(string(data), "\n", " ", -1)
	return strings.Replace(temp, "\t", " ", -1)
}
