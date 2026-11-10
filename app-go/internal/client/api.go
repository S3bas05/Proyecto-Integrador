package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"imageflow/app-go/internal/models"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) Health() error {
	request, err := http.NewRequest(http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("no se pudo conectar con el servidor: %w", err)
	}
	defer response.Body.Close()
	return checkResponse(response)
}

func (c *APIClient) Upload(filePath, category string) (models.Metadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return models.Metadata{}, fmt.Errorf("no se pudo abrir el archivo local: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return models.Metadata{}, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return models.Metadata{}, fmt.Errorf("no se pudo preparar el archivo: %w", err)
	}
	if err := writer.WriteField("category", category); err != nil {
		return models.Metadata{}, err
	}
	if err := writer.Close(); err != nil {
		return models.Metadata{}, err
	}

	request, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/images", &body)
	if err != nil {
		return models.Metadata{}, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := c.httpClient.Do(request)
	if err != nil {
		return models.Metadata{}, fmt.Errorf("falló el envío al servidor: %w", err)
	}
	defer response.Body.Close()
	if err := checkResponse(response); err != nil {
		return models.Metadata{}, err
	}

	var metadata models.Metadata
	if err := json.NewDecoder(response.Body).Decode(&metadata); err != nil {
		return models.Metadata{}, fmt.Errorf("no se pudo interpretar la respuesta del servidor: %w", err)
	}
	return metadata, nil
}

func (c *APIClient) ListImages(name, category string) ([]models.Metadata, error) {
	endpoint, err := url.Parse(c.baseURL + "/api/v1/images")
	if err != nil {
		return nil, err
	}
	query := endpoint.Query()
	if strings.TrimSpace(name) != "" {
		query.Set("name", name)
	}
	if strings.TrimSpace(category) != "" {
		query.Set("category", category)
	}
	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("no se pudo consultar el servidor: %w", err)
	}
	defer response.Body.Close()
	if err := checkResponse(response); err != nil {
		return nil, err
	}

	var images []models.Metadata
	if err := json.NewDecoder(response.Body).Decode(&images); err != nil {
		return nil, fmt.Errorf("no se pudo interpretar la lista de imágenes: %w", err)
	}
	return images, nil
}

func (c *APIClient) Download(id, destination string) error {
	endpoint := c.baseURL + "/api/v1/images/" + url.PathEscape(id) + "/download"
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("no se pudo descargar la imagen: %w", err)
	}
	defer response.Body.Close()
	if err := checkResponse(response); err != nil {
		return err
	}

	directory := filepath.Dir(destination)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("no se pudo preparar el directorio de descarga: %w", err)
	}
	tempFile, err := os.CreateTemp(directory, ".imageflow-download-*")
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo de descarga: %w", err)
	}
	tempName := tempFile.Name()
	defer os.Remove(tempName)

	if _, err := io.Copy(tempFile, response.Body); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("no se pudo guardar la descarga: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempName, destination); err != nil {
		return fmt.Errorf("no se pudo completar la descarga: %w", err)
	}
	return nil
}

func (c *APIClient) Delete(id string) error {
	endpoint := c.baseURL + "/api/v1/images/" + url.PathEscape(id)
	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("no se pudo enviar la solicitud de eliminación: %w", err)
	}
	defer response.Body.Close()
	return checkResponse(response)
}

func checkResponse(response *http.Response) error {
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	var payload map[string]string
	_ = json.NewDecoder(response.Body).Decode(&payload)
	if message := payload["error"]; message != "" {
		return fmt.Errorf("el servidor respondió %s: %s", response.Status, message)
	}
	return fmt.Errorf("el servidor respondió %s", response.Status)
}
