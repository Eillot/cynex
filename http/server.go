package http

import "net/http"

var server *Server

type Server string

func init() {
	server = new(Server)
}

func (s *Server) Run() {
	http.ListenAndServe(":8080", reactor)
}
