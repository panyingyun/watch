package loralog

import (
	"log"
	"os"
)

var debugLog = log.New(os.Stdout, "[Debug] ", log.LstdFlags)
var infoLog = log.New(os.Stdout, "[Info] ", log.LstdFlags)
var warnLog = log.New(os.Stdout, "[Warn] ", log.LstdFlags)
var errorLog = log.New(os.Stdout, "[Error] ", log.LstdFlags)

func Debugf(format string, v ...interface{}) {
	debugLog.Printf(format, v)
}

func Infof(format string, v ...interface{}) {
	infoLog.Printf(format, v)
}

func Warnf(format string, v ...interface{}) {
	warnLog.Printf(format, v)
}

func Errorf(format string, v ...interface{}) {
	errorLog.Printf(format, v)
}

func Debug(v ...interface{}) {
	debugLog.Println(v)
}

func Info(v ...interface{}) {
	infoLog.Println(v)
}

func Warn(v ...interface{}) {
	warnLog.Println(v)
}

func Error(v ...interface{}) {
	errorLog.Println(v)
}
