package controller

import (
	"fmt"

	"github.com/VladimirAzanza/url-shortener/internal/services"
	"github.com/gofiber/fiber/v2"
)

type FiberURLController struct {
	service *services.URLService
}

func NewFiberURLController(service *services.URLService) *FiberURLController {
	return &FiberURLController{service: service}
}

func (c *FiberURLController) HandlePost(ctx *fiber.Ctx) error {
	fmt.Println("Post handler")
	return nil
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "invalid request body", http.StatusBadRequest)
	// 	return
	// }
	// defer r.Body.Close()

	// originalURL := string(body)
	// shortURL := c.service.ShortenURL(originalURL)

	// w.Header().Set("Content-Type", "text/plain")
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(shortURL))
}

func (c *FiberURLController) HandleGet(ctx *fiber.Ctx) error {
	fmt.Println("Get handler")
	return nil
	// shortID := strings.TrimPrefix(r.URL.Path, "/")
	// if shortID == "" {
	// 	http.Error(w, "not ID addded", http.StatusBadRequest)
	// 	return
	// }

	// originalURL, exists := c.service.GetOriginalURL(shortID)
	// if !exists {
	// 	http.Error(w, "URL not found", http.StatusNotFound)
	// 	return
	// }

	// w.Header().Set("Location", originalURL)
	// w.WriteHeader(http.StatusTemporaryRedirect)
}
