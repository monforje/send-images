package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync"

	"send-images-backend/internal/logger"
	"send-images-backend/internal/util"
)

type Image struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Modified int64  `json:"modified"`
}

var imageExtPattern = regexp.MustCompile(`(?i)\.(jpe?g|png|gif|webp)$`)

func ImagesHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listImages(w, r, uploadDir)
		case http.MethodDelete:
			deleteImage(w, r, uploadDir)
		default:
			util.JSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

func listImages(w http.ResponseWriter, r *http.Request, dir string) {
	logger.Debug("Listing images...")

	entries, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("Cannot read dir: %v", err)
		util.JSONError(w, http.StatusInternalServerError, "Cannot read uploads directory")
		return
	}

	const maxWorkers = 16
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		images []Image
		sem    = make(chan struct{}, maxWorkers) // ограничим параллелизм
	)

	for _, entry := range entries {
		if entry.IsDir() || !imageExtPattern.MatchString(entry.Name()) {
			continue
		}

		name := entry.Name()
		fullPath := filepath.Join(dir, name)

		wg.Add(1)
		sem <- struct{}{}

		go func(name, path string) {
			defer wg.Done()
			defer func() { <-sem }()

			info, err := os.Stat(path)
			if err != nil {
				logger.Warn("Cannot stat file %s: %v", name, err)
				return
			}

			img := Image{
				Filename: name,
				URL:      "/uploads/" + name,
				Modified: info.ModTime().Unix(),
			}

			mu.Lock()
			images = append(images, img)
			mu.Unlock()
		}(name, fullPath)
	}

	wg.Wait()

	// Сортировка: новые выше
	sort.Slice(images, func(i, j int) bool {
		return images[i].Modified > images[j].Modified
	})

	paginated := paginateImages(images, r)
	util.JSON(w, http.StatusOK, map[string]interface{}{"images": paginated})
}

func paginateImages(images []Image, r *http.Request) []Image {
	limit, err1 := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, err2 := strconv.Atoi(r.URL.Query().Get("offset"))

	if err1 != nil || limit <= 0 {
		limit = 50 // дефолт
	}
	if err2 != nil || offset < 0 {
		offset = 0
	}

	if offset >= len(images) {
		return []Image{}
	}

	end := offset + limit
	if end > len(images) {
		end = len(images)
	}
	return images[offset:end]
}

func deleteImage(w http.ResponseWriter, r *http.Request, dir string) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		util.JSONError(w, http.StatusBadRequest, "Missing 'filename' parameter")
		return
	}

	safeName := filepath.Base(filename)
	fullPath := filepath.Join(dir, safeName)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		logger.Error("Failed to delete %s: %v", fullPath, err)
		util.JSONError(w, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	logger.Info("Deleted file: %s", safeName)
	util.JSON(w, http.StatusOK, map[string]string{
		"message":  "File deleted",
		"filename": safeName,
	})
}
