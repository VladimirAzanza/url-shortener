package main

import (
	"net/http"

	"github.com/VladimirAzanza/url-shortener/internal/app"
)

func main() {
	err := http.ListenAndServe(":8080", http.HandlerFunc(app.MainService))
	if err != nil {
		panic(err)
	}
}
