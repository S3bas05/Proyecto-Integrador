package models

import (
	"errors"
	"path/filepath"
	"strings"
	"time"
)

var ErrInvalidImage = errors.New("los datos de la imagen no son válidos")

// Image representa una imagen dentro del dominio de ImageFlow.
// Sus campos son privados para controlar cómo se crea y modifica el objeto.
type Image struct {
	id          string
	name        string
	storedName  string
	contentType string
	size        int64
	category    string
	createdAt   time.Time
}

// Metadata es la representación pública que se intercambia mediante JSON.
type Metadata struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

// Record contiene los metadatos y el nombre interno del archivo guardado.
// StoredName no se expone en las respuestas públicas de la API.
type Record struct {
	Metadata
	StoredName string `json:"stored_name"`
}

// NewImage es el constructor. Evita crear imágenes con información inválida.
func NewImage(id, name, storedName, contentType string, size int64, category string, createdAt time.Time) (Image, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	storedName = strings.TrimSpace(storedName)
	contentType = strings.TrimSpace(contentType)
	category = strings.TrimSpace(category)

	if id == "" || name == "" || storedName == "" || contentType == "" || size <= 0 {
		return Image{}, ErrInvalidImage
	}
	if len(name) > 120 || name == "." || name == ".." || strings.ContainsAny(name, `/\\`) {
		return Image{}, ErrInvalidImage
	}
	if filepath.Base(storedName) != storedName || strings.ContainsAny(storedName, `/\\`) {
		return Image{}, ErrInvalidImage
	}
	if category == "" {
		category = "sin-categoria"
	}
	if len(category) > 50 || strings.ContainsAny(category, "\r\n") {
		return Image{}, ErrInvalidImage
	}
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	return Image{
		id:          id,
		name:        name,
		storedName:  storedName,
		contentType: contentType,
		size:        size,
		category:    category,
		createdAt:   createdAt.UTC(),
	}, nil
}

// ImageFromRecord reconstruye el objeto de dominio desde los datos persistidos.
func ImageFromRecord(record Record) (Image, error) {
	return NewImage(
		record.ID,
		record.Name,
		record.StoredName,
		record.ContentType,
		record.Size,
		record.Category,
		record.CreatedAt,
	)
}

// Record convierte el objeto en una estructura lista para persistirse.
func (i Image) Record() Record {
	return Record{
		Metadata:   i.Metadata(),
		StoredName: i.storedName,
	}
}

// Metadata devuelve los datos que el cliente puede consultar.
func (i Image) Metadata() Metadata {
	return Metadata{
		ID:          i.id,
		Name:        i.name,
		ContentType: i.contentType,
		Size:        i.size,
		Category:    i.category,
		CreatedAt:   i.createdAt,
	}
}

func (i Image) ID() string           { return i.id }
func (i Image) Name() string         { return i.name }
func (i Image) StoredName() string   { return i.storedName }
func (i Image) ContentType() string  { return i.contentType }
func (i Image) Size() int64          { return i.size }
func (i Image) Category() string     { return i.category }
func (i Image) CreatedAt() time.Time { return i.createdAt }
