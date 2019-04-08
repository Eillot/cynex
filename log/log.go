package log

import (
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 文件输出等级，ERROR > WARNING > INFO > DEBUG
// 大于等于当前配置等级的日志，将被输出到文件
var Threshold = "INFO"

// 日志文件输出路径（文件夹地址）
var Dir = ""

// 日志输出配置初始化
const flag = log.Ldate | log.Ltime

var debug = log.New(os.Stdout, "DEBUG   ", flag)
var info = log.New(os.Stdout, "INFO    ", flag)
var waring = log.New(os.Stdout, "WARNING ", flag)
var error = log.New(os.Stdout, "ERROR   ", flag)

var debugF = log.New(logTargetFile(), "DEBUG   ", flag)
var infoF = log.New(logTargetFile(), "INFO    ", flag)
var warningF = log.New(logTargetFile(), "WARNING ", flag)
var errorF = log.New(logTargetFile(), "ERROR   ", flag)

// Debug output debug level logs
func Debug(v ...interface{}) {
	_, p, l, _ := runtime.Caller(1)
	wd, _ := os.Getwd()
	fileName := p[len(wd):] + ":" + strconv.Itoa(l)
	v2 := make([]interface{}, len(v)-1)
	v2 = append(v2, fileName)
	for _, vv := range v {
		v2 = append(v2, vv)
	}
	go debug.Println(v2...)
	if strings.ToLower(strings.TrimSpace(Threshold)) == "debug" {
		go debugF.Println(v2...)
	}
}

// Info output information level logs
func Info(v ...interface{}) {
	_, p, l, _ := runtime.Caller(1)
	wd, _ := os.Getwd()
	fileName := p[len(wd):] + ":" + strconv.Itoa(l)
	v2 := make([]interface{}, len(v)-1)
	v2 = append(v2, fileName)
	for _, vv := range v {
		v2 = append(v2, vv)
	}
	go info.Println(v2...)
	threshold := strings.ToLower(strings.TrimSpace(Threshold))
	if threshold == "debug" || threshold == "info" {
		go infoF.Println(v2...)
	}
}

// Warning output warning level logs
func Warning(v ...interface{}) {
	_, p, l, _ := runtime.Caller(1)
	wd, _ := os.Getwd()
	fileName := p[len(wd):] + ":" + strconv.Itoa(l)
	v2 := make([]interface{}, len(v)-1)
	v2 = append(v2, fileName)
	for _, vv := range v {
		v2 = append(v2, vv)
	}
	go waring.Println(v2...)
	threshold := strings.ToLower(strings.TrimSpace(Threshold))
	if threshold == "debug" || threshold == "info" || threshold == "warning" {
		go warningF.Println(v2...)
	}
}

// Error output error level logs
func Error(v ...interface{}) {
	_, p, l, _ := runtime.Caller(1)
	wd, _ := os.Getwd()
	fileName := p[len(wd):] + ":" + strconv.Itoa(l)
	v2 := make([]interface{}, len(v)-1)
	v2 = append(v2, fileName)
	for _, vv := range v {
		v2 = append(v2, vv)
	}
	go error.Println(v2...)
	go errorF.Println(v2...)
}

func logTargetFile() io.Writer {
	date := time.Now().Format("2006-01-02")
	var logFileName string
	if strings.TrimSpace(Dir) == "" {
		logFileName = "./" + date + ".log"
	} else {
		logFileName = Dir + "/" + date + ".log"
	}
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return logFile
}
