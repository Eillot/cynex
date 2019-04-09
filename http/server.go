package http

import (
	"cynex/log"
	"net/http"
)

var Server *server

type server struct {
	handler     *reactor
	downloadDir string
}

func init() {
	Server = &server{
		handler:     handler,
		downloadDir: ".",
	}
}

func (s *server) Start() {
	log.Info("正在启动服务...")
	http.ListenAndServe(":8080", s.handler)
}

func StopServer() {

}
