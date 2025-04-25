package logger

import (
	"log"
	"os"
	"strings"
)

var (
	debugMode bool
)

// init determines whether debug mode is enabled by checking the DEBUG
// environment variable. If set to "true", debug messages will be printed.
func init() {
	debugMode = strings.ToLower(os.Getenv("DEBUG")) == "true"
}

// Debug logs a debug message formatted with the provided arguments,
// but only if debugMode is enabled. The message is prefixed with "[DEBUG]".

func Debug(msg string, args ...interface{}) {
	if debugMode {
		log.Printf("[DEBUG] "+msg, args...)
	}
}

// Info logs an information message formatted with the provided arguments.
// The message is prefixed with "[INFO]".
func Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

// Error logs an error message formatted with the provided arguments.
// The message is prefixed
func Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}
