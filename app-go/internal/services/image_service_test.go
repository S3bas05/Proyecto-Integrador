package services_test

import (
	"errors"
	"testing"

	"imageflow/app-go/internal/services"
	"imageflow/app-go/internal/storage"
)

func newService(t *testing.T) *services.ImageService {
	t.Helper()
	store, err := storage.NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}
	return services.NewImageService(store)
}

func TestUploadRejectsOversizedImage(t *testing.T) {
	service := newService(t)
	data := make([]byte, int(services.MaxImageSize)+1)

	_, err := service.Upload("imagen.png", "prueba", data)
	if !errors.Is(err, services.ErrImageTooLarge) {
		t.Fatalf("Upload() error = %v, want %v", err, services.ErrImageTooLarge)
	}
}

func TestUploadRejectsPathInFileName(t *testing.T) {
	service := newService(t)
	pngData := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89,
	}

	_, err := service.Upload("../imagen.png", "prueba", pngData)
	if !errors.Is(err, services.ErrInvalidInput) {
		t.Fatalf("Upload() error = %v, want %v", err, services.ErrInvalidInput)
	}
}
