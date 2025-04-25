package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	debugMode = strings.ToLower(os.Getenv("DEBUG")) == "true"
	logChan   = make(chan string, 1000)
	stopOnce  sync.Once
	wg        sync.WaitGroup

	logOut io.Writer = os.Stdout
	logger           = log.New(logOut, "", 0)
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[34m"
)

func init() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range logChan {
			logger.Println(msg)
		}
	}()
}

// Shutdown завершает логгер (ждёт окончания вывода)
func Shutdown() {
	stopOnce.Do(func() {
		close(logChan)
		wg.Wait()
	})
}

// SetOutput позволяет направить лог в файл или буфер
func SetOutput(w io.Writer) {
	logOut = w
	logger = log.New(logOut, "", 0)
}

func logf(color, level, emoji, msg string, args ...interface{}) {
	formatted := format(level, emoji, msg, args...)
	sendLog(color + formatted + reset)
}

func format(level, emoji, msg string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	full := fmt.Sprintf(msg, args...)
	return fmt.Sprintf("%s [%s] %s %s", timestamp, level, emoji, full)
}

func Debug(msg string, args ...interface{}) {
	if debugMode {
		logf(blue, "DEBUG", "🐛", msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	logf(green, "INFO", "ℹ️", msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logf(yellow, "WARN", "⚠️", msg, args...)
}

func Error(msg string, args ...interface{}) {
	logf(red, "ERROR", "❌", msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logf(red, "FATAL", "💀", msg, args...)
	time.Sleep(100 * time.Millisecond)
	os.Exit(1)
}

func sendLog(msg string) {
	select {
	case logChan <- msg:
	default:
		logger.Println("[LOGGER] Log channel full — writing directly:")
		logger.Println(msg)
	}
}

func IsDebug() bool {
	return debugMode
}
