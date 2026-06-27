package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/PedroScheurer/product-service/internal/apperrors"
	"github.com/PedroScheurer/product-service/internal/storage"
)

type MediaController struct {
	storage    *storage.MediaStorage
	maxVideoMB int64
	maxImageMB int64
}

func NewMediaController(store *storage.MediaStorage, maxVideoMB, maxImageMB int64) *MediaController {
	return &MediaController{
		storage:    store,
		maxVideoMB: maxVideoMB,
		maxImageMB: maxImageMB,
	}
}

func (c *MediaController) RegisterRoutes(r chi.Router) {
	r.Post("/ws/products/upload", c.uploadMedia)
	r.Handle("/media/*", c.serveMedia())
}

func (c *MediaController) uploadMedia(w http.ResponseWriter, r *http.Request) {
	headers, err := parseRequiredHeaders(r)
	if err != nil {
		writeWsError(w, err)
		return
	}
	if headers.userType != adminType {
		writeWsError(w, apperrors.NewAuthenticationError("Usuário sem Permissão!"))
		return
	}

	kind := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("kind")))
	maxBytes := c.maxImageMB * 1024 * 1024
	if kind == storage.KindVideo {
		maxBytes = c.maxVideoMB * 1024 * 1024
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	if err := r.ParseMultipartForm(maxBytes); err != nil {
		writeWsError(w, fmt.Errorf("arquivo muito grande ou inválido"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeWsError(w, fmt.Errorf("campo 'file' é obrigatório"))
		return
	}
	defer file.Close()

	publicPath, err := c.storage.Save(kind, header.Filename, header.Header.Get("Content-Type"), file)
	if err != nil {
		writeWsError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"url": publicPath})
}

func (c *MediaController) serveMedia() http.Handler {
	return http.StripPrefix("/media/", http.FileServer(http.Dir(c.storage.RootDir())))
}
