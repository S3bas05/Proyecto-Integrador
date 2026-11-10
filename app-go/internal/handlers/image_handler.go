package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"

	"imageflow/app-go/internal/services"
)

type ImageHandler struct {
	service *services.ImageService
}

func NewImageHandler(service *services.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

func (h *ImageHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/images", h.upload)
	mux.HandleFunc("GET /api/v1/images", h.list)
	mux.HandleFunc("GET /api/v1/images/{id}/download", h.download)
	mux.HandleFunc("DELETE /api/v1/images/{id}", h.delete)
}

func (h *ImageHandler) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, services.MaxImageSize+(1<<20))
	parseErr := r.ParseMultipartForm(services.MaxImageSize)
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	if parseErr != nil {
		var maxErr *http.MaxBytesError
		if errors.As(parseErr, &maxErr) || errors.Is(parseErr, multipart.ErrMessageTooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, services.ErrImageTooLarge.Error())
			return
		}
		writeError(w, http.StatusBadRequest, "solicitud multipart inválida; debe enviar el campo file")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "no se recibió el archivo en el campo file")
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, services.MaxImageSize+1))
	if err != nil {
		writeError(w, http.StatusBadRequest, "no se pudo leer el archivo enviado")
		return
	}
	metadata, err := h.service.Upload(header.Filename, r.FormValue("category"), content)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, metadata)
}

func (h *ImageHandler) list(w http.ResponseWriter, r *http.Request) {
	images, err := h.service.List(r.URL.Query().Get("name"), r.URL.Query().Get("category"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, images)
}

func (h *ImageHandler) download(w http.ResponseWriter, r *http.Request) {
	metadata, content, err := h.service.Download(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", metadata.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	if disposition := mime.FormatMediaType("attachment", map[string]string{"filename": metadata.Name}); disposition != "" {
		w.Header().Set("Content-Disposition", disposition)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func (h *ImageHandler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "imagen eliminada correctamente"})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrUnsupportedFormat):
		writeError(w, http.StatusUnsupportedMediaType, err.Error())
	case errors.Is(err, services.ErrImageTooLarge):
		writeError(w, http.StatusRequestEntityTooLarge, err.Error())
	case errors.Is(err, services.ErrImageNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "ocurrió un error interno al procesar la operación")
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
