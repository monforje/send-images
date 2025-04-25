package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"send-images-backend/internal/logger"
)

// JSON отправляет JSON-ответ с указанным статусом и логирует.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Буфер для безопасного логирования
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(true) // безопаснее для браузера

	if err := encoder.Encode(payload); err != nil {
		logger.Error("Failed to encode JSON: %v", err)
		return
	}

	// Пишем в response
	if _, err := w.Write(buf.Bytes()); err != nil {
		logger.Error("Failed to write JSON response: %v", err)
		return
	}

	if logger.IsDebug() {
		logger.Debug("Response status %d, payload: %s", status, buf.String())
	}
}

// JSONError — единый стиль ошибок с логированием.
func JSONError(w http.ResponseWriter, status int, message string) error {
	errObj := map[string]string{
		"error": message,
		"code":  http.StatusText(status),
	}
	JSON(w, status, errObj)

	return fmt.Errorf("http %d: %s", status, message)
}
