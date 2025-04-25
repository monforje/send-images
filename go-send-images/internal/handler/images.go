package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"send-images-backend/internal/logger"
	"send-images-backend/internal/util"
)

func ImagesHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listImages(w, r, uploadDir)
		case http.MethodDelete:
			deleteImage(w, r, uploadDir)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func listImages(w http.ResponseWriter, r *http.Request, dir string) {
	logger.Debug("Listing images...")

	entries, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("Cannot read dir: %v", err)
		http.Error(w, "Cannot read uploads dir", http.StatusInternalServerError)
		return
	}

	var images []map[string]string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if ok, _ := regexp.MatchString(`(?i)\.(jpe?g|png|gif|webp|bmp)$`, name); !ok {
			continue
		}
		images = append(images, map[string]string{
			"filename": name,
			"url":      "/uploads/" + name,
		})
	}

	// сортировка по имени
	sort.Slice(images, func(i, j int) bool {
		return images[i]["filename"] < images[j]["filename"]
	})

	// пагинация
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = len(images)
	}
	if offset > len(images) {
		offset = len(images)
	}
	end := offset + limit
	if end > len(images) {
		end = len(images)
	}
	paged := images[offset:end]

	util.JSON(w, http.StatusOK, map[string]interface{}{"images": paged})
}

func deleteImage(w http.ResponseWriter, r *http.Request, dir string) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Missing filename parameter", http.StatusBadRequest)
		return
	}

	safeName := filepath.Base(filename)
	fullPath := filepath.Join(dir, safeName)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
		} else {
			logger.Error("Failed to delete file: %v", err)
			http.Error(w, "Cannot delete file", http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Deleted file: %s", fullPath)
	util.JSON(w, http.StatusOK, map[string]string{
		"message":  "File deleted",
		"filename": filename,
	})
}
