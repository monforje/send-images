package main

import (
	"log"
	"net/http"
	"os"

	"send-images-backend/internal/handler"
	"send-images-backend/internal/logger"
)

const (
	uploadDir = "./uploads"
	addr      = "0.0.0.0:9999"
)

// main starts the server. It creates the upload directory if it doesn't exist,
// sets up the HTTP handlers and starts listening on the specified address.
func main() {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload dir: %v", err)
	}

	logger.Info("Upload dir ready at %s", uploadDir)

	mux := http.NewServeMux()
	mux.HandleFunc("/images", handler.ImagesHandler(uploadDir))
	mux.HandleFunc("/upload", handler.UploadHandler(uploadDir))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	logger.Info("Server running at %s", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(h http.Handler) http.Handler {
	// withCORS wraps an HTTP handler to add CORS headers to the response,
	// allowing requests from any origin. It supports GET, POST, DELETE, and
	// OPTIONS methods, and allows the "Content-Type" header. For OPTIONS
	// requests, it responds with a 204 status code without further processing.

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
