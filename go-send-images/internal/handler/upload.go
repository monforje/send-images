package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"send-images-backend/internal/db"
	"send-images-backend/internal/logger"
	"send-images-backend/internal/model"
	"send-images-backend/internal/util"
)

const maxFileSize = 5 << 20 // 5MB

var (
	allowedTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/gif":  true,
	}
	filenameSanitizer = regexp.MustCompile(`[^\w\-]`)
	bufPool           = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 32*1024)
			return &b
		},
	}
)

func UploadHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			util.JSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		logger.Debug("Received upload from %s", r.RemoteAddr)

		r.Body = http.MaxBytesReader(w, r.Body, 20*maxFileSize+1024)
		if err := r.ParseMultipartForm(20 * maxFileSize); err != nil {
			logger.Error("Malformed form or too large: %v", err)
			util.JSONError(w, http.StatusBadRequest, "Request too large or malformed")
			return
		}

		files := r.MultipartForm.File["file"]
		if len(files) == 0 {
			util.JSONError(w, http.StatusBadRequest, "No files uploaded")
			return
		}

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			logger.Error("Upload dir creation failed: %v", err)
			util.JSONError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		var (
			wg            sync.WaitGroup
			mu            sync.Mutex
			uploadedFiles = make([]map[string]string, 0, len(files))
		)

		for _, fh := range files {
			fh := fh
			wg.Add(1)
			go func(header *multipart.FileHeader) {
				defer wg.Done()

				entry, err := handleUpload(header, uploadDir, r)
				if err != nil {
					logger.Warn("Upload failed: %v", err)
					return
				}
				mu.Lock()
				uploadedFiles = append(uploadedFiles, entry)
				mu.Unlock()
			}(fh)
		}

		wg.Wait()

		if len(uploadedFiles) == 0 {
			util.JSONError(w, http.StatusBadRequest, "No valid images uploaded")
			return
		}

		util.JSON(w, http.StatusOK, map[string]interface{}{
			"message": "Files uploaded successfully",
			"files":   uploadedFiles,
		})
	}
}

func setupCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleUpload(header *multipart.FileHeader, uploadDir string, r *http.Request) (map[string]string, error) {
	src, err := header.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	if !isAllowedImage(src) {
		return nil, fmt.Errorf("unsupported file type: %s", header.Filename)
	}
	src.Seek(0, io.SeekStart)

	filename := generateFilename(header.Filename)
	dstPath := filepath.Join(uploadDir, filename)

	if err := saveFile(dstPath, src); err != nil {
		return nil, fmt.Errorf("saving failed: %w", err)
	}

	url := "/uploads/" + filename
	modified := time.Now().Unix()

	image := model.Image{
		Filename: filename,
		URL:      url,
		Modified: modified,
	}

	_, err = db.Collection("images").InsertOne(r.Context(), image)
	if err != nil {
		logger.Error("Failed to save image metadata to Mongo: %v", err)
	}

	logger.Info("Saved file: %s", dstPath)

	return map[string]string{
		"filename": filename,
		"url":      url,
	}, nil
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
	safe := filenameSanitizer.ReplaceAllString(base, "_")
	return fmt.Sprintf("%s-%s%s", uuid.New().String(), safe, ext)
}

func saveFile(dstPath string, src io.Reader) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)

	_, err = io.CopyBuffer(dst, src, *bufPtr)
	return err
}
