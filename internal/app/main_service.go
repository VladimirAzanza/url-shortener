package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var urlStore = make(map[int]string, 0)

func MainService(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		handleGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var originalURL string
	err := json.NewDecoder(r.Body).Decode(&originalURL)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortID := 123
	_, exists := urlStore[shortID]

	if exists {
		shortID += 10
		urlStore[shortID] = originalURL
	} else {
		urlStore[shortID] = originalURL
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%d", (shortID))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(shortURL))

}

func handleGet(w http.ResponseWriter, r *http.Request) {

}
