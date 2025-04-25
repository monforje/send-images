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

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[34m"
)

func format(level, msg string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	full := fmt.Sprintf(msg, args...)
	return fmt.Sprintf("[%s] [%s] %s", timestamp, level, full)
}

func Debug(msg string, args ...interface{}) {
	if debugMode {
		logger.Println(blue + format("DEBUG", msg, args...) + reset)
	}
}

func Info(msg string, args ...interface{}) {
	logger.Println(green + format("INFO", msg, args...) + reset)
}

func Warn(msg string, args ...interface{}) {
	logger.Println(yellow + format("WARN", msg, args...) + reset)
}

func Error(msg string, args ...interface{}) {
	logger.Println(red + format("ERROR", msg, args...) + reset)
}

func Fatal(msg string, args ...interface{}) {
	logger.Fatalln(red + format("FATAL", msg, args...) + reset)
}
