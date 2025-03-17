package controller

import (
	"io"
	"net/http"
	"strings"

	"github.com/VladimirAzanza/url-shortener/internal/services"
)

type URLController struct {
	service *services.URLService
}

func NewURLController(service *services.URLService) *URLController {
	return &URLController{service: service}
}

func (c *URLController) HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)
	shortURL := c.service.ShortenURL(originalURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(shortURL))
}

func (c *URLController) HandleGet(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/")
	if shortID == "" {
		http.Error(w, "not ID addded", http.StatusBadRequest)
		return
	}

	originalURL, exists := c.service.GetOriginalURL(shortID)
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
