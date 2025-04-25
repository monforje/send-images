package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"send-images-backend/internal/logger"
	"send-images-backend/internal/util"

	"github.com/google/uuid"
)

const maxFileSize = 5 << 20 // 5MB

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

func UploadHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS для фронта
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		logger.Debug("Receiving upload request from %s", r.RemoteAddr)

		// Ограничение размера тела запроса
		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize*20+512) // увеличено для нескольких файлов

		if err := r.ParseMultipartForm(maxFileSize * 20); err != nil {
			logger.Error("Request too large or malformed: %v", err)
			http.Error(w, "Request too large", http.StatusBadRequest)
			return
		}

		form := r.MultipartForm
		files := form.File["file"]

		if len(files) == 0 {
			http.Error(w, "No files uploaded", http.StatusBadRequest)
			return
		}

		var uploadedFiles []map[string]string

		for _, header := range files {
			src, err := header.Open()
			if err != nil {
				logger.Error("Failed to open uploaded file: %v", err)
				continue
			}

			if !isAllowedImage(src) {
				logger.Warn("Skipping unsupported file type: %s", header.Filename)
				src.Close()
				continue
			}

			src.Seek(0, io.SeekStart)

			filename := generateFilename(header.Filename)
			dstPath := filepath.Join(uploadDir, filename)

			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				logger.Error("Failed to create upload dir: %v", err)
				src.Close()
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			if err := saveFile(dstPath, src); err != nil {
				logger.Error("Failed to save file: %v", err)
				src.Close()
				continue
			}

			src.Close()

			logger.Info("Saved file: %s", dstPath)

			uploadedFiles = append(uploadedFiles, map[string]string{
				"filename": filename,
				"url":      "/uploads/" + filename,
			})
		}

		if len(uploadedFiles) == 0 {
			http.Error(w, "No valid images uploaded", http.StatusBadRequest)
			return
		}

		util.JSON(w, http.StatusOK, map[string]interface{}{
			"message": "Files uploaded successfully",
			"files":   uploadedFiles,
		})
	}
}

func isAllowedImage(file io.ReadSeeker) bool {
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	file.Seek(0, io.SeekStart)
	return allowedTypes[contentType]
}

func generateFilename(original string) string {
	ext := filepath.Ext(original)
	base := strings.TrimSuffix(original, ext)
	safe := regexp.MustCompile(`[^\w\-]`).ReplaceAllString(base, "_")
	return fmt.Sprintf("%s-%s%s", uuid.New().String(), safe, ext)
}

func saveFile(path string, src io.Reader) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}
