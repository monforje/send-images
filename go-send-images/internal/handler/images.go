package handler

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"send-images-backend/internal/db"
	"send-images-backend/internal/logger"
	"send-images-backend/internal/model"
	"send-images-backend/internal/util"
)

func ImagesHandler(uploadDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listImagesFromMongo(w, r)
		case http.MethodDelete:
			deleteImageFromDiskAndMongo(w, r, uploadDir)
		default:
			util.JSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

// listImagesFromMongo читает список изображений из MongoDB
func listImagesFromMongo(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	findOptions.SetSort(bson.D{{Key: "modified", Value: -1}})

	cur, err := db.Collection("images").Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		logger.Error("Failed to query images from Mongo: %v", err)
		util.JSONError(w, http.StatusInternalServerError, "Database error")
		return
	}
	defer cur.Close(context.Background())

	var images []model.Image
	if err := cur.All(context.Background(), &images); err != nil {
		logger.Error("Failed to decode images: %v", err)
		util.JSONError(w, http.StatusInternalServerError, "Decoding error")
		return
	}

	util.JSON(w, http.StatusOK, map[string]interface{}{
		"images": images,
	})
}

// deleteImageFromDiskAndMongo удаляет файл и метаданные из MongoDB
func deleteImageFromDiskAndMongo(w http.ResponseWriter, r *http.Request, uploadDir string) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		util.JSONError(w, http.StatusBadRequest, "Missing 'filename' parameter")
		return
	}

	safeName := filepath.Base(filename) // защита от path traversal
	fullPath := filepath.Join(uploadDir, safeName)

	// Удаляем файл
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		logger.Error("Failed to delete file %s: %v", fullPath, err)
		util.JSONError(w, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	// Удаляем запись из Mongo
	res, err := db.Collection("images").DeleteOne(context.Background(), bson.M{"filename": safeName})
	if err != nil {
		logger.Error("Failed to delete metadata from Mongo: %v", err)
	}

	logger.Info("Deleted file and metadata: %s", safeName)
	util.JSON(w, http.StatusOK, map[string]interface{}{
		"message":  "File deleted",
		"filename": safeName,
		"deleted":  res.DeletedCount,
	})
}
