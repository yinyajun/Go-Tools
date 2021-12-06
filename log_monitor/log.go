package log_monitor

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func InitLogger(level string) {
	flag := log.Ldate | log.Ltime | log.Lshortfile
	DebugLogger = log.New(os.Stdout, "[Debug]", flag)
	InfoLogger = log.New(os.Stdout, "[Info]", flag)
	WarningLogger = log.New(os.Stdout, "[Warning]", flag)
	ErrorLogger = log.New(os.Stdout, "[Error]", flag)

	switch level {
	case "info":
		DebugLogger = log.New(ioutil.Discard, "", flag)
	case "warning":
		DebugLogger = log.New(ioutil.Discard, "", flag)
		InfoLogger = log.New(ioutil.Discard, "", flag)
	case "error":
		DebugLogger = log.New(ioutil.Discard, "", flag)
		InfoLogger = log.New(ioutil.Discard, "", flag)
		WarningLogger = log.New(ioutil.Discard, "", flag)
	}
}
