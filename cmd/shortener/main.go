package main

import (
	"net/http"

	"github.com/VladimirAzanza/url-shortener/internal/controller"
)

func main() {
	err := http.ListenAndServe(":8080", http.HandlerFunc(controller.MainService))
	if err != nil {
		panic(err)
	}
}
