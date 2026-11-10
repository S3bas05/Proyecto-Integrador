package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"imageflow/app-go/internal/models"
)

var ErrNotFound = errors.New("imagen no encontrada")
var ErrAlreadyExists = errors.New("ya existe una imagen con ese identificador")

// LocalStorage separa los archivos binarios de sus metadatos JSON.
type LocalStorage struct {
	dataDir      string
	imagesDir    string
	metadataPath string
	items        map[string]models.Record
	mu           sync.RWMutex
}

func NewLocalStorage(dataDir string) (*LocalStorage, error) {
	if strings.TrimSpace(dataDir) == "" {
		return nil, errors.New("la ruta de almacenamiento está vacía")
	}

	absoluteDir, err := filepath.Abs(dataDir)
	if err != nil {
		return nil, fmt.Errorf("no se pudo resolver la ruta de datos: %w", err)
	}

	imagesDir := filepath.Join(absoluteDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return nil, fmt.Errorf("no se pudo crear el directorio de imágenes: %w", err)
	}

	store := &LocalStorage{
		dataDir:      absoluteDir,
		imagesDir:    imagesDir,
		metadataPath: filepath.Join(absoluteDir, "metadata.json"),
		items:        make(map[string]models.Record),
	}

	content, err := os.ReadFile(store.metadataPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("no se pudieron leer los metadatos: %w", err)
	}
	if err == nil && len(content) > 0 {
		var records []models.Record
		if err := json.Unmarshal(content, &records); err != nil {
			return nil, fmt.Errorf("el archivo metadata.json no contiene JSON válido: %w", err)
		}
		for _, record := range records {
			if _, exists := store.items[record.ID]; exists {
				return nil, fmt.Errorf("metadata.json contiene un ID duplicado: %s", record.ID)
			}
			if _, err := models.ImageFromRecord(record); err != nil {
				return nil, fmt.Errorf("se encontró un registro inválido en metadata.json: %w", err)
			}
			if !safeStoredName(record.StoredName) {
				return nil, errors.New("metadata.json contiene un nombre interno de archivo no válido")
			}
			store.items[record.ID] = record
		}
	}
	if errors.Is(err, os.ErrNotExist) {
		if err := store.persistLocked(); err != nil {
			return nil, err
		}
	}

	return store, nil
}

func (s *LocalStorage) Save(record models.Record, data []byte) error {
	if !safeStoredName(record.StoredName) || record.ID == "" || len(data) == 0 {
		return errors.New("registro o archivo no válido")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[record.ID]; exists {
		return ErrAlreadyExists
	}

	targetPath := filepath.Join(s.imagesDir, record.StoredName)
	tempPath := targetPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return fmt.Errorf("no se pudo escribir la imagen: %w", err)
	}
	if err := os.Rename(tempPath, targetPath); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("no se pudo guardar la imagen: %w", err)
	}

	s.items[record.ID] = record
	if err := s.persistLocked(); err != nil {
		delete(s.items, record.ID)
		_ = os.Remove(targetPath)
		return fmt.Errorf("no se pudieron guardar los metadatos: %w", err)
	}
	return nil
}

func (s *LocalStorage) List() []models.Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := make([]models.Record, 0, len(s.items))
	for _, record := range s.items {
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.After(records[j].CreatedAt)
	})
	return records
}

func (s *LocalStorage) Get(id string) (models.Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.items[id]
	if !exists {
		return models.Record{}, ErrNotFound
	}
	return record, nil
}

func (s *LocalStorage) ReadFile(id string) (models.Record, []byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, exists := s.items[id]
	if !exists {
		return models.Record{}, nil, ErrNotFound
	}
	content, err := os.ReadFile(filepath.Join(s.imagesDir, record.StoredName))
	if err != nil {
		return models.Record{}, nil, fmt.Errorf("no se pudo leer el archivo de imagen: %w", err)
	}
	return record, content, nil
}

func (s *LocalStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, exists := s.items[id]
	if !exists {
		return ErrNotFound
	}

	filePath := filepath.Join(s.imagesDir, record.StoredName)
	backup, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("no se pudo leer la imagen antes de eliminarla: %w", err)
	}
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("no se pudo eliminar el archivo: %w", err)
	}

	delete(s.items, id)
	if err := s.persistLocked(); err != nil {
		s.items[id] = record
		_ = os.WriteFile(filePath, backup, 0600)
		return fmt.Errorf("no se pudieron actualizar los metadatos: %w", err)
	}
	return nil
}

func (s *LocalStorage) persistLocked() error {
	records := make([]models.Record, 0, len(s.items))
	for _, record := range s.items {
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })

	content, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("no se pudieron convertir los metadatos a JSON: %w", err)
	}
	tempPath := s.metadataPath + ".tmp"
	if err := os.WriteFile(tempPath, content, 0600); err != nil {
		return fmt.Errorf("no se pudo escribir el archivo temporal de metadatos: %w", err)
	}
	if err := os.Rename(tempPath, s.metadataPath); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("no se pudo actualizar metadata.json: %w", err)
	}
	return nil
}

func safeStoredName(name string) bool {
	return name != "" && name != "." && name != ".." && filepath.Base(name) == name && !strings.ContainsAny(name, `/\\`)
}
