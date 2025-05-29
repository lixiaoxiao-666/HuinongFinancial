package utils

import (
	"log"
	"os"
)

var (
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// LogError 记录错误日志
func LogError(message string, err error) {
	if err != nil {
		errorLogger.Printf("%s: %v", message, err)
	} else {
		errorLogger.Printf("%s", message)
	}
}

// LogInfo 记录信息日志
func LogInfo(message string) {
	infoLogger.Printf("%s", message)
}

// LogWarn 记录警告日志
func LogWarn(message string) {
	errorLogger.Printf("WARN: %s", message)
}
