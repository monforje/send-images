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
	uploadDir     = getEnv("UPLOAD_DIR", "./uploads")
	addr          = getEnv("ADDR", "0.0.0.0:9999")
	allowedOrigin = getEnv("ALLOWED_ORIGIN", "http://localhost:3000")
)

func main() {
	// Проверим и создадим директорию
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Fatal("Failed to create upload dir: %v", err)
	}
	logger.Info("Upload dir ready at %s", uploadDir)

	// Роуты
	mux := http.NewServeMux()
	mux.HandleFunc("/images", handler.ImagesHandler(uploadDir))
	mux.HandleFunc("/upload", handler.UploadHandler(uploadDir))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))
	mux.HandleFunc("/health", healthHandler)

	// HTTP сервер
	server := &http.Server{
		Addr:    addr,
		Handler: withCORS(mux),
	}

	// Остановка по Ctrl+C или SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("🚀 Server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error: %v", err)
		}
	}()

	<-stop
	logger.Info("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Завершаем сервер
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Shutdown error: %v", err)
	}

	// Завершаем логгер (если буферизация)
	logger.Shutdown()

	logger.Info("✅ Server gracefully stopped")
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	// Можно заменить на настоящую проверку
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
