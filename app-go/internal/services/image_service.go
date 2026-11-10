package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"imageflow/app-go/internal/models"
	"imageflow/app-go/internal/storage"
)

const MaxImageSize int64 = 10 * 1024 * 1024 // 10 MiB

var (
	ErrInvalidInput      = errors.New("los datos enviados no son válidos")
	ErrUnsupportedFormat = errors.New("formato de imagen no permitido; use JPG, PNG, GIF o WEBP")
	ErrImageTooLarge     = errors.New("la imagen supera el límite de 10 MB")
	ErrImageNotFound     = errors.New("imagen no encontrada")
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

type ImageService struct {
	storage *storage.LocalStorage
}

func NewImageService(store *storage.LocalStorage) *ImageService {
	return &ImageService{storage: store}
}

func (s *ImageService) Upload(fileName, category string, data []byte) (models.Metadata, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || len(data) == 0 || strings.ContainsAny(fileName, `/\\`) || fileName == "." || fileName == ".." {
		return models.Metadata{}, ErrInvalidInput
	}
	if len(data) > int(MaxImageSize) {
		return models.Metadata{}, ErrImageTooLarge
	}

	contentType := http.DetectContentType(data[:min(len(data), 512)])
	extension, allowed := allowedImageTypes[contentType]
	if !allowed {
		return models.Metadata{}, ErrUnsupportedFormat
	}

	id, err := newID()
	if err != nil {
		return models.Metadata{}, fmt.Errorf("no se pudo generar el identificador: %w", err)
	}

	// La extensión se obtiene del contenido real, no del nombre enviado por el cliente.
	storedName := id + extension
	image, err := models.NewImage(id, fileName, storedName, contentType, int64(len(data)), category, time.Now().UTC())
	if err != nil {
		return models.Metadata{}, ErrInvalidInput
	}
	if err := s.storage.Save(image.Record(), data); err != nil {
		return models.Metadata{}, fmt.Errorf("no se pudo almacenar la imagen: %w", err)
	}
	return image.Metadata(), nil
}

func (s *ImageService) List(name, category string) ([]models.Metadata, error) {
	records := s.storage.List()
	name = strings.ToLower(strings.TrimSpace(name))
	category = strings.ToLower(strings.TrimSpace(category))

	result := make([]models.Metadata, 0, len(records))
	for _, record := range records {
		image, err := models.ImageFromRecord(record)
		if err != nil {
			return nil, fmt.Errorf("metadatos inválidos: %w", err)
		}
		if name != "" && !strings.Contains(strings.ToLower(image.Name()), name) {
			continue
		}
		if category != "" && !strings.Contains(strings.ToLower(image.Category()), category) {
			continue
		}
		result = append(result, image.Metadata())
	}
	return result, nil
}

func (s *ImageService) Download(id string) (models.Metadata, []byte, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return models.Metadata{}, nil, ErrInvalidInput
	}
	record, content, err := s.storage.ReadFile(id)
	if errors.Is(err, storage.ErrNotFound) {
		return models.Metadata{}, nil, ErrImageNotFound
	}
	if err != nil {
		return models.Metadata{}, nil, err
	}
	image, err := models.ImageFromRecord(record)
	if err != nil {
		return models.Metadata{}, nil, err
	}
	return image.Metadata(), content, nil
}

func (s *ImageService) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrInvalidInput
	}
	err := s.storage.Delete(id)
	if errors.Is(err, storage.ErrNotFound) {
		return ErrImageNotFound
	}
	return err
}

func newID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
