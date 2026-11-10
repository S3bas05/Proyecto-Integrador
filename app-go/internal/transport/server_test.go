package transport_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	apiClient "imageflow/app-go/internal/client"
	"imageflow/app-go/internal/services"
	"imageflow/app-go/internal/storage"
	"imageflow/app-go/internal/transport"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	store, err := storage.NewLocalStorage(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}
	service := services.NewImageService(store)
	return httptest.NewServer(transport.NewServer(service))
}

func TestHealthEndpoint(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	response, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health error = %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	var payload map[string]string
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("JSON decode error = %v", err)
	}
	if payload["status"] != "ok" {
		t.Fatalf("status field = %q, want %q", payload["status"], "ok")
	}
}

func TestClientServerImageOperations(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()
	api := apiClient.NewAPIClient(server.URL)

	// PNG mínimo de prueba. Se crea como archivo local porque el cliente lo carga desde disco.
	pngData := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0x60, 0x00, 0x02, 0x00,
		0x00, 0x05, 0x00, 0x01, 0xa5, 0xf6, 0x45, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
		0xae, 0x42, 0x60, 0x82,
	}
	localFile := filepath.Join(t.TempDir(), "prueba.png")
	if err := os.WriteFile(localFile, pngData, 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	metadata, err := api.Upload(localFile, "paisaje")
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if metadata.ID == "" || metadata.Category != "paisaje" {
		t.Fatalf("metadata inesperada: %+v", metadata)
	}

	images, err := api.ListImages("prueba", "paisaje")
	if err != nil {
		t.Fatalf("ListImages() error = %v", err)
	}
	if len(images) != 1 || images[0].ID != metadata.ID {
		t.Fatalf("ListImages() = %+v, se esperaba una imagen", images)
	}

	downloadPath := filepath.Join(t.TempDir(), "descarga.png")
	if err := api.Download(metadata.ID, downloadPath); err != nil {
		t.Fatalf("Download() error = %v", err)
	}
	downloaded, err := os.ReadFile(downloadPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !bytes.Equal(downloaded, pngData) {
		t.Fatal("el archivo descargado no coincide con el original")
	}

	if err := api.Delete(metadata.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	images, err = api.ListImages("", "")
	if err != nil {
		t.Fatalf("ListImages() después de borrar: %v", err)
	}
	if len(images) != 0 {
		t.Fatalf("quedaron %d imágenes después de eliminar", len(images))
	}
}

func TestUploadRejectsUnsupportedFile(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "not-image.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("esto no es una imagen"))
	_ = writer.WriteField("category", "prueba")
	_ = writer.Close()

	request, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/images", &body)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusUnsupportedMediaType)
	}
}

func TestHealthRejectsUnsupportedMethod(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	request, err := http.NewRequest(http.MethodPost, server.URL+"/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestDownloadMissingImageReturnsNotFound(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	response, err := http.Get(server.URL + "/api/v1/images/no-existe/download")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusNotFound)
	}
}
