package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"send-images-backend/internal/handler"
	"send-images-backend/internal/logger"
	"syscall"
	"time"
)

var (
	uploadDir = getEnv("UPLOAD_DIR", "./uploads")
	addr      = getEnv("ADDR", "0.0.0.0:9999")
)

func main() {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Error("Failed to create upload dir: %v", err)
		os.Exit(1)
	}
	logger.Info("Upload dir ready at %s", uploadDir)

	mux := http.NewServeMux()
	mux.HandleFunc("/images", handler.ImagesHandler(uploadDir))
	mux.HandleFunc("/upload", handler.UploadHandler(uploadDir))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	server := &http.Server{
		Addr:    addr,
		Handler: withCORS(mux),
	}

	go func() {
		logger.Info("Server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error: %v", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Shutdown error: %v", err)
	}
}

func withCORS(h http.Handler) http.Handler {
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

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
