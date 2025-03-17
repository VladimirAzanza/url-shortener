package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var urlStore = make(map[string]string, 0)

func MainService(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		handleGet(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)
	shortID := "123"
	_, exists := urlStore[shortID]

	if exists {
		shortID += "10"
		urlStore[shortID] = originalURL
	} else {
		urlStore[shortID] = originalURL
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%s", (shortID))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(shortURL))

}

func handleGet(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/")
	if shortID == "" {
		http.Error(w, "not ID addded", http.StatusBadRequest)
		return
	}

	originalURL, exists := urlStore[shortID]
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
