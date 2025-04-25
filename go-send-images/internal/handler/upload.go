package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"send-images-backend/internal/logger"
	"send-images-backend/internal/util"
)

const maxFileSize = 5 << 20 // 5MB

func UploadHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		logger.Debug("Receiving upload request from %s", r.RemoteAddr)

		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize+512)
		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			logger.Error("File too large: %v", err)
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Error("Missing file field: %v", err)
			http.Error(w, "Missing file field", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if !isValidImage(file) {
			http.Error(w, "Only images allowed", http.StatusBadRequest)
			return
		}

		file.Seek(0, io.SeekStart)

		filename := generateFilename(header.Filename)
		dstPath := filepath.Join(uploadDir, filename)

		if err := saveFile(dstPath, file); err != nil {
			logger.Error("Failed to write file: %v", err)
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

func isValidImage(file io.ReadSeeker) bool {
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	contentType := http.DetectContentType(buf[:n])
	file.Seek(0, io.SeekStart)
	return strings.HasPrefix(contentType, "image/")
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
