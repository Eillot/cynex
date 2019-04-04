package log

import (
	"log"
	"os"
)

var info = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
var waring = log.New(os.Stdout, "WARNING ", log.Ldate|log.Ltime|log.Lshortfile)
var error = log.New(os.Stdout, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

func Info(v ...interface{}) {
	info.Println(v)
}

func Warning(v ...interface{}) {
	waring.Println(v)
}

func Error(v ...interface{}) {
	error.Println(v)
}
