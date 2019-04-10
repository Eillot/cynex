package react

import (
	"cynex/log"
	"net/http"
)

var Server *server

type server struct {
	handler     *handler
	downloadDir string
}

func init() {

	Server = &server{
		handler:     defaultHandler,
		downloadDir: ".",
	}
}

func (s *server) Start() {
	log.Info("正在启动服务...")
	http.ListenAndServe(":8080", s.handler)
}
