package transport

import (
	"encoding/json"
	"net/http"

	"imageflow/app-go/internal/handlers"
	"imageflow/app-go/internal/services"
)

func NewServer(imageService *services.ImageService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	handlers.NewImageHandler(imageService).Register(mux)
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"service": "ImageFlow API",
		"status":  "ok",
	})
}
