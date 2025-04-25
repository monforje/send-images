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
	logChan   = make(chan string, 1000) // буферизированный канал
	stopOnce  sync.Once
	wg        sync.WaitGroup

	logOut = io.MultiWriter(os.Stdout) // Можно добавить os.OpenFile(...) для записи в файл
	logger = log.New(logOut, "", 0)
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

func Shutdown() {
	stopOnce.Do(func() {
		close(logChan)
		wg.Wait()
	})
}

func format(level, emoji, msg string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	full := fmt.Sprintf(msg, args...)
	return fmt.Sprintf("%s [%s] %s %s%s", timestamp, level, emoji, full, reset)
}

func Debug(msg string, args ...interface{}) {
	if debugMode {
		sendLog(blue + format("DEBUG", "🐛", msg, args...))
	}
}

func Info(msg string, args ...interface{}) {
	sendLog(green + format("INFO", "ℹ️", msg, args...))
}

func Warn(msg string, args ...interface{}) {
	sendLog(yellow + format("WARN", "⚠️", msg, args...))
}

func Error(msg string, args ...interface{}) {
	sendLog(red + format("ERROR", "❌", msg, args...))
}

func Fatal(msg string, args ...interface{}) {
	sendLog(red + format("FATAL", "💀", msg, args...))
	time.Sleep(100 * time.Millisecond) // на случай async-канала
	os.Exit(1)
}

func sendLog(msg string) {
	select {
	case logChan <- msg:
	default:
		// fallback при переполнении
		logger.Println("[LOGGER] Log channel full, writing directly:")
		logger.Println(msg)
	}
}

func IsDebug() bool {
	return debugMode
}
