package router

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

func (data *CommRouter) dealWithResponse() {
	if getContentType() != ContentTypeNon {
		data.content = getContentType()
	} else {
		data.content = data.GetContentType()
	}

	// check function error
	if data.Err != nil {
		if len(data.RspBody) == 0 {
			if data.ErrCode == nil {
				data.ErrCode = "FAIL"
			}
			if data.Message == "" {
				data.Message = data.Err.Error()
			}
			if data.content == ContentTypeJson {
				data.RspBody = []byte(fmt.Sprintf(JsonErrResponseStr, data.ErrCode, data.Err, data.Message, data.uid))
				data.SetHeader("Content-Type", "application/json;charset=utf-8")
			} else if data.content == ContentTypeXml {
				data.RspBody = []byte(fmt.Sprintf(XmlErrResponseStr, data.ErrCode, data.Message, data.uid))
				data.SetHeader("Content-Type", "application/xml;charset=utf-8")
			} else {
				data.RspBody = []byte(data.Message)
			}
		}

		data.status = http.StatusAccepted
		return
	} else if len(data.RspBody) > 0 {
		data.status = http.StatusOK
		return
	}

	switch data.content {
	case ContentTypeJson:
		if data.Result != nil {
			// convert data.Result struct to byte
			if data.RspBody, data.Err = json.Marshal(data.Result); data.Err == nil {
				data.status = http.StatusOK
			} else {
				data.RspBody = []byte(fmt.Sprintf(JsonErrResponseStr, "FAIL", data.Err, "Response convert to json wrong", data.uid))
			}
		} else if len(data.RspMap) > 0 {
			// convert data.Result struct to byte
			if data.RspBody, data.Err = json.Marshal(data.RspMap); data.Err == nil {
				data.status = http.StatusOK
			} else {
				data.RspBody = []byte(fmt.Sprintf(JsonErrResponseStr, "FAIL", data.Err, "Response convert to json wrong", data.uid))
			}
		}
		data.SetHeader("Content-Type", "application/json;charset=utf-8")
	case ContentTypeXml:
		if data.Result != nil {
			// convert data.Result struct to byte
			if data.RspBody, data.Err = xml.Marshal(data.Result); data.Err == nil {
				data.status = http.StatusOK
			} else {
				data.RspBody = []byte(fmt.Sprintf(XmlErrResponseStr, "FAIL", data.Err, data.uid))
			}
		} else if len(data.RspMap) > 0 {
			// convert data.RspMap
			if data.RspBody, data.Err = xml.Marshal(XmlMap(data.RspMap)); data.Err == nil {
				data.status = http.StatusOK
			} else {
				data.RspBody = []byte(fmt.Sprintf(XmlErrResponseStr, "FAIL", data.Err, data.uid))
			}
		}
		data.SetHeader("Content-Type", "application/xml;charset=utf-8")
	}

	if data.status == 0 {
		data.status = http.StatusAccepted
	}
	return
}
