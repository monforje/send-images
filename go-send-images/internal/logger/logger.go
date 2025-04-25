package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	debugMode = strings.ToLower(os.Getenv("DEBUG")) == "true"
	logger    = log.New(os.Stdout, "", 0)
)

func format(level, msg string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, fmt.Sprintf(msg, args...))
}

func Debug(msg string, args ...interface{}) {
	if debugMode {
		logger.Println(format("DEBUG", msg, args...))
	}
}

func Info(msg string, args ...interface{}) {
	logger.Println(format("INFO", msg, args...))
}

func Warn(msg string, args ...interface{}) {
	logger.Println(format("WARN", msg, args...))
}

func Error(msg string, args ...interface{}) {
	logger.Println(format("ERROR", msg, args...))
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatalln(format("FATAL", msg, args...))
}
