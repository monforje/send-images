package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"send-images-backend/internal/logger"
	"send-images-backend/internal/util"
)

// ImagesHandler returns an HTTP handler for the /images endpoint. It supports two
// methods: GET to list all images in the upload directory, and DELETE to delete an
// image by filename. The upload directory is specified as an argument to the
// function. The handler logs debug information about the requests it receives and
// returns to the client.
func ImagesHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listImages(w, uploadDir)
		case http.MethodDelete:
			deleteImage(w, r, uploadDir)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// listImages writes a JSON response containing a list of image files in the specified directory.
// It reads the directory contents, filters out non-image files, and constructs a list of maps
// with each containing a filename and a URL. If the directory cannot be read, an error response
// is sent to the client.

func listImages(w http.ResponseWriter, dir string) {
	logger.Debug("Listing images...")

	entries, err := os.ReadDir(dir)
	if err != nil {
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

	util.JSON(w, http.StatusOK, map[string]interface{}{"images": images})
}

// deleteImage deletes an image file from the specified directory based on the filename query parameter.
// If the file does not exist, it returns a 404 status code. If the deletion fails, it returns a 500 status code.
// Otherwise, it returns a 204 status code.
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
			http.Error(w, "Cannot delete file", http.StatusInternalServerError)
		}
		return
	}

	logger.Info("Deleted file: %s", fullPath)
	w.WriteHeader(http.StatusNoContent)
}
