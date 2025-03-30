package controller

import (
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
	baseURL := ctx.BaseURL()
	originalURL := ctx.BodyRaw()
	shortID := c.service.ShortenURL(string(originalURL))

	return ctx.Status(fiber.StatusOK).SendString(baseURL + "/" + shortID)
}

func (c *FiberURLController) HandleGet(ctx *fiber.Ctx) error {
	shortID := ctx.Params("id")

	originalURL, exists := c.service.GetOriginalURL(shortID)
	if !exists {
		return ctx.Status(fiber.StatusNotFound).SendString("URL not found")
	}
	return ctx.Redirect(originalURL, fiber.StatusTemporaryRedirect)
}
