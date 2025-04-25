package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"send-images-backend/internal/logger"
	"send-images-backend/internal/util"
)

const maxFileSize = 5 << 20 // 5MB

// UploadHandler returns an HTTP handler for uploading images. It accepts POST requests
// with a "file" field containing the image. The handler saves the image to the
// specified upload directory, and returns a JSON response with the saved filename
// and the URL of the file. The handler also logs the request and the saved file.
func UploadHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		logger.Debug("Receiving upload request from %s", r.RemoteAddr)

		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize+512)
		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Missing file field", http.StatusBadRequest)
			return
		}
		defer file.Close()

		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		contentType := http.DetectContentType(buf[:n])
		if !strings.HasPrefix(contentType, "image/") {
			http.Error(w, "Only images allowed", http.StatusBadRequest)
			return
		}
		file.Seek(0, io.SeekStart)

		logger.Debug("File accepted: %s, type: %s", header.Filename, contentType)

		ext := filepath.Ext(header.Filename)
		base := strings.TrimSuffix(header.Filename, ext)
		safe := regexp.MustCompile(`[^\w\-]`).ReplaceAllString(base, "_")
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		filename := fmt.Sprintf("%s-%s%s", timestamp, safe, ext)

		dstPath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			return
		}

		logger.Info("Saved file: %s", dstPath)

		util.JSON(w, http.StatusOK, map[string]interface{}{
			"message":  "Файл загружен",
			"filename": filename,
			"url":      "/uploads/" + filename,
		})
	}
}
