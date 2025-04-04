package controller

import (
	"github.com/VladimirAzanza/url-shortener/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type FiberURLController struct {
	service *services.URLService
}

func NewFiberURLController(service *services.URLService) *FiberURLController {
	return &FiberURLController{service: service}
}

// HandlePost Post a URL to shorten
// @Summary Shorten a URL
// @Description Create a short URL from the original URL
// @Tags URLs
// @Accept plain
// @Produce plain
// @Param originalUrl body string true "Original URL to be shortened"
// @Success 201 {string} string "Returns the shortened URL"
// @Router / [post]
func (c *FiberURLController) HandlePost(ctx *fiber.Ctx) error {
	baseURL := ctx.BaseURL()
	originalURL := ctx.BodyRaw()
	shortID := c.service.ShortenURL(string(originalURL))

	return ctx.Status(fiber.StatusCreated).SendString(baseURL + "/" + shortID)
}

// HandleGet godoc
// @Summary Redirect to original URL
// @Description Redirects to the original URL using the short ID
// @Tags URLs
// @Produce plain
// @Param id path string true "Short URL ID"
// @Success 307 "Redirects to original URL"
// @Failure 404 {string} string "Not found if short ID doesn't exist"
// @Router /{id} [get]
func (c *FiberURLController) HandleGet(ctx *fiber.Ctx) error {
	shortID := ctx.Params("id")

	originalURL, exists := c.service.GetOriginalURL(shortID)
	if !exists {
		return ctx.Status(fiber.StatusNotFound).SendString("URL not found")
	}
	log.Info().Str("shortID", shortID).Str("originalURL", originalURL).Msg("Redirect to original URL")
	return ctx.Redirect(originalURL, fiber.StatusTemporaryRedirect)
}
