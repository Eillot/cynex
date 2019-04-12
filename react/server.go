package react

import (
	"cynex/conf"
	"cynex/log"
	"net/http"
	"strings"
)

var Server *server

type server struct {
	handler     *handler
	downloadDir string
}

const (
	wd      string = ""
	confDir string = "/conf"
)

func init() {

	Server = &server{
		handler:     defaultHandler,
		downloadDir: ".",
	}
}

func (s *server) Start() {
	// 配置加载
	log.Info("正在加载配置...")
	configs := make(map[string]string)
	configsWd, err := conf.Load(wd)
	if err != nil {
		log.Error("配置加载错误：===>", err)
	}
	for key, val := range configsWd {
		if strings.TrimSpace(key) != "" && strings.TrimSpace(val) != "" {
			configs[key] = val
		}
	}
	configsDir, err := conf.Load(confDir)
	if err != nil {
		log.Error("配置加载错误：===>", err)
	}
	for key, val := range configsDir {
		if strings.TrimSpace(key) != "" && strings.TrimSpace(val) != "" {
			configs[key] = val
		}
	}
	//
	log.Info("正在启动服务...")
	http.ListenAndServe(":8080", s.handler)
}
