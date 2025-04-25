package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"send-images-backend/internal/logger"
)

// JSON отправляет JSON-ответ с указанным статусом и логирует (если debug включен)
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(true) // предотвращает XSS в HTML

	if err := encoder.Encode(payload); err != nil {
		logger.Error("❌ Failed to encode JSON: %v", err)
		return
	}

	if _, err := w.Write(buf.Bytes()); err != nil {
		logger.Error("❌ Failed to write JSON response: %v", err)
		return
	}

	if logger.IsDebug() {
		logger.Debug("✅ JSON Response [%d]: %s", status, buf.String())
	}
}

// JSONError возвращает ошибку с заданным статусом и логирует её
func JSONError(w http.ResponseWriter, status int, message string) error {
	errObj := map[string]interface{}{
		"error":   message,
		"code":    http.StatusText(status),
		"status":  status,
		"success": false,
	}
	JSON(w, status, errObj)
	return fmt.Errorf("http %d: %s", status, message)
}
