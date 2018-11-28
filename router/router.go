package router

import (
	"net/http"
	"strings"
)

type RouterPath struct {
	method   string
	path     string
	function func(*CommRouter)
}

type Server struct {
	content ContentType
	server  *http.Server
}

var defaultContentType = ContentTypeNon
var routers map[string]*RouterPath

func GetFunction(path string) *RouterPath {
	return routers[path]
}

func SetDefaultContentType(contentType string) {
	switch strings.ToUpper(contentType) {
	case "json":
		defaultContentType = ContentTypeJson
	case "xml":
		defaultContentType = ContentTypeXml
	case "byte":
		defaultContentType = ContentTypeByte
	}
}

func NewRouter(server *http.Server) *Server {
	routers = make(map[string]*RouterPath)
	return &Server{server: server}
}

func (s *Server) Post(path string, function func(*CommRouter)) {
	routers[path] = &RouterPath{
		method:   "POST",
		path:     path,
		function: function,
	}
}

func (s *Server) Get(path string, function func(*CommRouter)) {
	routers[path] = &RouterPath{
		method:   "GET",
		path:     path,
		function: function,
	}
}

func (s *Server) Put(path string, function func(*CommRouter)) {
	routers[path] = &RouterPath{
		method:   "PUT",
		path:     path,
		function: function,
	}
}

func (s *Server) Delete(path string, function func(*CommRouter)) {
	routers[path] = &RouterPath{
		method:   "DELETE",
		path:     path,
		function: function,
	}
}

func (s *Server) Run() {
	s.server.Handler = http.HandlerFunc(CommonDealRequest)

	if err := s.server.ListenAndServe(); err != nil {
		panic(err)
	}
}
