package storage_test

import (
	"bytes"
	"testing"
	"time"

	"imageflow/app-go/internal/models"
	"imageflow/app-go/internal/storage"
)

func TestLocalStoragePersistsImageAndMetadata(t *testing.T) {
	directory := t.TempDir()
	store, err := storage.NewLocalStorage(directory)
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}

	image, err := models.NewImage("test-id", "foto.png", "test-id.png", "image/png", 4, "prueba", time.Now())
	if err != nil {
		t.Fatalf("NewImage() error = %v", err)
	}
	original := []byte{1, 2, 3, 4}
	if err := store.Save(image.Record(), original); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	restartedStore, err := storage.NewLocalStorage(directory)
	if err != nil {
		t.Fatalf("NewLocalStorage() after restart error = %v", err)
	}
	record, downloaded, err := restartedStore.ReadFile("test-id")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if record.Name != "foto.png" {
		t.Fatalf("record.Name = %q, want %q", record.Name, "foto.png")
	}
	if !bytes.Equal(downloaded, original) {
		t.Fatalf("ReadFile() = %v, want %v", downloaded, original)
	}
}
