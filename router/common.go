package router

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContentType int

var (
	ContentTypeNon  ContentType = 0 // default non
	ContentTypeJson ContentType = 1 // application/json
	ContentTypeXml  ContentType = 2 // application/xml, text/xml
	ContentTypeByte ContentType = 3 //
	ContentTypeText ContentType = 4 // text/plain
	ContentTypeHtml ContentType = 5 // text/html
	ContentTypeData ContentType = 6 // multipart/form-data

	ErrContentTypeUnknown = errors.New("Content-Type Not Known")
)

const (
	XmlErrResponseStr = `<xml><code>%v</code><message>%v</message><uid>%s</uid></xml>`

	JsonErrResponseStr = `{"error":{"code":"%v", "error":"%v", "message":"%v"}, "uid":"%s"}`
)

// make request event id
// UID = IP + Time + Rand  (28 bytes)
func MakeUid(r *http.Request, timestamp time.Time) string {
	ip := strings.Split(strings.Split(r.RemoteAddr, ":")[0], ".")

	rand.Seed(timestamp.UnixNano())
	var bte = make([]byte, 8)
	for i := 0; i < 8; i++ {
		bte[i] = byte(rand.Intn(26) + 97)
	}

	return fmt.Sprintf("%s%s%s",
		byteToString(ip),
		timestamp.Format("060102150405"),
		string(bte),
	)
}

// can convert 0-255 to string
func byteToString(data []string) string {
	var result = make([]byte, 0)
	for _, row := range data {
		d, _ := strconv.Atoi(row)
		result = append(result, byte(d/26+48), byte(d%26+97))
	}
	return string(result)
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
