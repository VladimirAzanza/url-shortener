package controller

import (
	"fmt"

	"github.com/VladimirAzanza/url-shortener/internal/constants"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
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

// HandleAPIPost Post a URL to shorten
// @Summary Shorten a URL
// @Description Create a short URL from the original URL
// @Tags URLs
// @Accept plain
// @Produce plain
// @Param request body dto.ShortenRequestDTO true "Original URL to be shortened"
// @Success 201 {object} dto.ShortenResponseDTO "Returns the shortened URL"
// @Failure 500 {object} map[string]string "Failed to parse request"
// @Router / [post]
func (c *FiberURLController) HandleAPIPost(ctx *fiber.Ctx) error {
	var shortenRequestDTO dto.ShortenRequestDTO
	if err := ctx.BodyParser(&shortenRequestDTO); err != nil {
		log.Err(err).Msg(constants.MsgFailedToParseBody)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.MsgFailedToParseBody,
		})
	}

	shortID := c.service.ShortenAPIURL(ctx.Context(), &shortenRequestDTO)

	fullURL := fmt.Sprintf("%s/%s", ctx.BaseURL(), shortID)
	response := dto.ShortenResponseDTO{
		Result: fullURL,
	}

	// jsonBytes, err := json.Marshal(response)
	// if err != nil {
	// 	log.Err(err).Msg("Failed to marshal response")
	// 	return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	// }

	// log.Info().Msg("Successfully shortened the url, shortID: " + response.Result)
	// ctx.Status(fiber.StatusCreated)
	// ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	// return ctx.Send(jsonBytes)

	log.Info().Msg("Successfully shortened the url, shortID" + response.Result)
	return ctx.Status(fiber.StatusCreated).JSON(response)
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
