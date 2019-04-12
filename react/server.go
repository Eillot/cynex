package react

import (
	"cynex/conf"
	"cynex/log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var Server *server

type server struct {
	handler     *handler
	downloadDir string
}

const (
	defaultHttpPort  string = "80"
	defaultHttpsPort string = "443"
)

const (
	// 配置文件路径
	wd      string = ""
	confDir string = "/conf"
	// 配置文件标识
	httpPort      string = "http.port"
	httpsEnable   string = "https.enable"
	httpsPort     string = "https.port"
	httpsKeyFile  string = "https.key_dir"
	httpsCertFile string = "https.cert_dir"
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
	// 使用配置启动服务
	log.Info("正在启动服务...")
	if s.getHttpsEnable(configs) {
		go http.ListenAndServeTLS(":"+s.getHttpsPort(configs), s.getHttpsCert(configs), s.getHttpsKey(configs), s.handler)
	}
	http.ListenAndServe(":"+s.getHttpPort(configs), s.handler)
}

func (s *server) getHttpPort(configs map[string]string) string {
	if configs[httpPort] != "" {
		return configs[httpPort]
	}
	return defaultHttpPort
}

func (s *server) getHttpsEnable(configs map[string]string) bool {
	if configs[httpsEnable] == "" {
		return false
	}
	reg, _ := regexp.Compile("[0-9]")
	if reg.MatchString(configs[httpsEnable]) {
		n, err := strconv.Atoi(configs[httpsEnable])
		if err != nil {
			log.Warning("配置项：https.enable 不符合语法规范")
			return false
		}
		if n > 0 {
			return true
		}
		return false
	}
	var ons = []string{"true", "on", "open", "start"}
	for _, on := range ons {
		if on == strings.ToLower(configs[httpsEnable]) {
			return true
		}
	}
	return false
}

func (s *server) getHttpsPort(configs map[string]string) string {
	if configs[httpsPort] != "" {
		return configs[httpsPort]
	}
	return defaultHttpsPort
}

func (s *server) getHttpsCert(configs map[string]string) string {
	dir := configs[httpsCertFile]
	if strings.Index(dir, "/") != 0 {
		if strings.Index(dir, "./") != 0 {
			wd, _ := os.Getwd()
			dir = wd + "/" + dir
		}
	}
	return dir
}

func (s *server) getHttpsKey(configs map[string]string) string {
	dir := configs[httpsKeyFile]
	if strings.Index(dir, "/") != 0 {
		if strings.Index(dir, "./") != 0 {
			wd, _ := os.Getwd()
			dir = wd + "/" + dir
		}
	}
}
