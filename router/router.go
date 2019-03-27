package router

import (
	"net/http"
	"strings"
)

// Set glabol server
var glabol *Server

const defaultPath = "DEFAULT"

type Server struct {
	maxRead     int64
	content     ContentType
	server      *http.Server
	contentType ContentType
	routers     map[string]*RouterPath
}

type RouterPath struct {
	method   string
	path     string
	function []func(*CommRouter)
}

// make new router
func NewRouter(server *http.Server) *Server {
	return &Server{
		maxRead:     1024 * 1024,
		server:      server,
		contentType: ContentTypeNon,
		routers:     make(map[string]*RouterPath),
	}
}

// when use response convert, cat default content type
// json, xml, byte
func (s *Server) SetDefaultContentType(contentType string) {
	switch strings.ToLower(contentType) {
	case "json":
		s.contentType = ContentTypeJson
	case "xml":
		s.contentType = ContentTypeXml
	case "byte":
		s.contentType = ContentTypeByte
	default:
		panic("Set Default Content-Type Wrong")
	}
}

// max requesty body read, default: 1024*1024 bytes
func (s *Server) SetMaxBodyRead(number int) {
	if number < 512 {
		panic("Min body lenght is 512 bytes")
	}
	s.maxRead = int64(number)
}

func (s *Server) Post(path string, functions ...func(*CommRouter)) {
	s.routers[path] = &RouterPath{
		method:   "POST",
		path:     path,
		function: functions,
	}
}

func (s *Server) Get(path string, functions ...func(*CommRouter)) {
	s.routers[path] = &RouterPath{
		method:   "GET",
		path:     path,
		function: functions,
	}
}

func (s *Server) Put(path string, functions ...func(*CommRouter)) {
	s.routers[path] = &RouterPath{
		method:   "PUT",
		path:     path,
		function: functions,
	}
}

func (s *Server) Delete(path string, functions ...func(*CommRouter)) {
	s.routers[path] = &RouterPath{
		method:   "DELETE",
		path:     path,
		function: functions,
	}
}

func (s *Server) Default(functions ...func(*CommRouter)) {
	s.routers[defaultPath] = &RouterPath{
		function: functions,
	}
}

func (s *Server) Run() {
	glabol = s
	glabol.server.Handler = http.HandlerFunc(CommonDealWith)

	if err := glabol.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func getContentType() ContentType {
	return glabol.contentType
}

func getFunction(path string) (*RouterPath, bool) {
	f, ok := glabol.routers[path]
	return f, ok
}

func getMaxBodyRead() int64 {
	return glabol.maxRead
}
