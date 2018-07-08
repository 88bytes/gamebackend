package utils

import (
	"log"
	"os"
)

//ILogger represents  the log interface
type ILogger interface {
	Println(v ...interface{})
	Fatal(v ...interface{})
}

var (
	// Logger 是gameServer用来打日志的
	Logger ILogger = log.New(os.Stderr, "", log.LstdFlags|log.Llongfile)
)
