package router

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
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

	s.SetDefaultContentType("xml")

	s.Post("/sandbox", CheckModel, Sandbox)
	s.Run()
}

func CheckModel(data *CommRouter) {
	data.SetResponseConvert()

	result, err := data.GetRequestMap()
	if err != nil {
		data.PutError("Can't convert to map")
	}
	if v, ok := result["mid"]; !ok {
		data.PutError("Request not have mid value")
	} else if _, err := strconv.Atoi(v); err != nil {
		data.PutError("Mid is not", "number")
		// data.Err = errors.New("Mid is not number")
	}
}

func Sandbox(data *CommRouter) {
	req := &SandboxRequest{}
	body, _ := data.GetBody()

	if data.GetContentType() == ContentTypeJson {
		if err := json.Unmarshal(body, req); err != nil {
			data.Err = err
			data.Message = "Sandbox json unmarshal error: " + err.Error()
			return
		}
	} else if data.GetContentType() == ContentTypeXml {
		if err := xml.Unmarshal(body, req); err != nil {
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
		data.PutError("Model Wrong")
	}
	return
}
