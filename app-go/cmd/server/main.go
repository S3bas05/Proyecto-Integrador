package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"imageflow/app-go/internal/services"
	"imageflow/app-go/internal/storage"
	"imageflow/app-go/internal/transport"
)

func main() {
	host := envOrDefault("IMAGEFLOW_HOST", "0.0.0.0")
	port := envOrDefault("IMAGEFLOW_PORT", "8080")
	dataDir := envOrDefault("IMAGEFLOW_DATA_DIR", "./data")

	imageStorage, err := storage.NewLocalStorage(dataDir)
	if err != nil {
		log.Fatalf("No se pudo preparar el almacenamiento: %v", err)
	}

	imageService := services.NewImageService(imageStorage)
	server := &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           transport.NewServer(imageService),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("ImageFlow Server iniciado en %s:%s", host, port)
	log.Printf("Directorio de datos: %s", dataDir)
	log.Println("Endpoint de prueba: GET /health")
	log.Println("API de imágenes: /api/v1/images")

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
