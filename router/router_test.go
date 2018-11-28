package router

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"testing"
)

type SandboxRequest struct {
	Model int    `xml:"model" json:"model"`
	Mid   string `xml:"mid"  json:"mid"`
}

type SandboxResponse struct {
	Message string `xml:"message" json:"message"`
}

func TestSandbox(t *testing.T) {
	s := NewRouter(&http.Server{
		Addr: ":9080",
	})
	s.Post("/sandbox", Sandbox)
	s.Run()
}

func Sandbox(data *CommRouter) {
	req := &SandboxRequest{}
	if data.ContentType == ContentTypeJson {
		if err := json.Unmarshal(data.ReqBody, req); err != nil {
			data.Err = err
			data.Message = "Sandbox xml unmarshal error: " + err.Error()
			return
		}
	} else if data.ContentType == ContentTypeXml {
		if err := xml.Unmarshal(data.ReqBody, req); err != nil {
			data.Err = err
			data.Message = "Sandbox xml unmarshal error: " + err.Error()
			return
		}
	}

	if req.Model == 1 {
		data.Result = &SandboxResponse{Message: "Result SUCCESS"}
	} else if req.Model == 2 {
		data.RspMap = make(map[string]string)
		data.RspMap["message"] = "Map SUCCESS"
	} else {
		data.RspBody = []byte("Byte SUCCESS")
	}
	return
}
