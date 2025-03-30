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
	shortID := ctx.Params("id")

	originalURL, exists := c.service.GetOriginalURL(shortID)
	if !exists {
		return ctx.Status(fiber.StatusNotFound).SendString("URL not found")
	}
	return ctx.Redirect(originalURL, fiber.StatusTemporaryRedirect)
}
