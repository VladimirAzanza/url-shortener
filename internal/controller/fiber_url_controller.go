package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/VladimirAzanza/url-shortener/internal/constants"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
	"github.com/VladimirAzanza/url-shortener/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type FiberURLController struct {
	service services.IURLService
}

func NewFiberURLController(service services.IURLService) *FiberURLController {
	return &FiberURLController{
		service: service,
	}
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
	shortID, err := c.service.ShortenURL(ctx.UserContext(), string(originalURL))
	if err != nil {
		log.Error().Err(err).Msg("Error at shorten api url")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).SendString(baseURL + "/" + shortID)
}

// @Summary Verifies the connection to the DB
// @Description Ping to the DB
// @Tags DB
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Success"
// @Failure 500 {object} map[string]string "Can not connect to the Database"
// @Router /ping [get]
func (c *FiberURLController) GetDBPing(ctx *fiber.Ctx) error {
	pingCtx, cancel := context.WithTimeout(ctx.Context(), 5*time.Second)
	defer cancel()

	if err := c.service.PingDB(pingCtx); err != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Can not connect to the Database",
			"error":   err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"storage_type": c.service.GetStorageType(),
	})
}

// HandleAPIPost Post a URL to shorten
// @Summary Shorten a URL
// @Description Create a short URL from the original URL
// @Tags API
// @Accept plain
// @Produce plain
// @Param request body dto.ShortenRequestDTO true "Original URL to be shortened"
// @Success 201 {object} dto.ShortenResponseDTO "Returns the shortened URL"
// @Failure 500 {object} map[string]string "When internal server error occurs"
// @Router /api/shorten [post]
func (c *FiberURLController) HandleAPIPost(ctx *fiber.Ctx) error {
	var shortenRequestDTO dto.ShortenRequestDTO
	if err := ctx.BodyParser(&shortenRequestDTO); err != nil {
		log.Err(err).Msg(constants.MsgFailedToParseBody)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.MsgFailedToParseBody,
		})
	}

	shortID, err := c.service.ShortenAPIURL(ctx.UserContext(), &shortenRequestDTO)
	if err != nil {
		log.Error().Err(err).Msg("Error at shorten api url")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fullURL := fmt.Sprintf("%s/%s", ctx.BaseURL(), shortID)
	response := dto.ShortenResponseDTO{
		Result: fullURL,
	}

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
// @Failure 408 {string} string "Request timeout"
// @Router /{id} [get]
func (c *FiberURLController) HandleGet(ctx *fiber.Ctx) error {
	shortID := ctx.Params("id")

	reqCtx, cancel := context.WithTimeout(ctx.UserContext(), 1*time.Second)
	defer cancel()

	originalURL, exists := c.service.GetOriginalURL(reqCtx, shortID)
	switch err := reqCtx.Err(); {
	case errors.Is(err, context.DeadlineExceeded):
		log.Warn().Str("shortID", shortID).Msg("Request timeout exceeded (server-side)")
		return ctx.Status(fiber.StatusRequestTimeout).SendString("Request timeout")
	case errors.Is(err, context.Canceled):
		log.Warn().Str("shortID", shortID).Msg("Request canceled by client")
		return ctx.Status(499).SendString("Client closed connection")
	case err != nil:
		log.Error().Err(err).Msg("Unexpected Error")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	if !exists {
		return ctx.Status(fiber.StatusNotFound).SendString("URL not found")
	}

	log.Info().Str("shortID", shortID).Str("originalURL", originalURL).Msg("Redirect to original URL")
	return ctx.Redirect(originalURL, fiber.StatusTemporaryRedirect)
}

// HandleAPIPostBatch Shorten multiple URLs in batch
// @Summary Shorten multiple URLs in a single request
// @Description Accepts a batch of URLs and returns their shortened versions
// @Tags API
// @Accept json
// @Produce json
// @Param request body []dto.BatchRequestDTO true "Array of URLs to shorten"
// @Success 201 {array} dto.BatchResponseDTO "Returns an array of shortened URLs"
// @Failure 400 {object} map[string]string "When request body is invalid or empty"
// @Failure 500 {object} map[string]string "When internal server error occurs"
// @Router /api/shorten/batch [post]
func (c *FiberURLController) HandleAPIPostBatch(ctx *fiber.Ctx) error {
	var batchRequestDTO []dto.BatchRequestDTO
	if err := ctx.BodyParser(&batchRequestDTO); err != nil {
		log.Err(err).Msg(constants.MsgFailedToParseBody)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.MsgFailedToParseBody,
		})
	}

	if len(batchRequestDTO) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty batch",
		})
	}

	responses := make([]dto.BatchResponseDTO, 0, len(batchRequestDTO))
	for _, req := range batchRequestDTO {
		shortID, err := c.service.BatchShortenURL(ctx.UserContext(), req)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		responses = append(responses, dto.BatchResponseDTO{
			CorrelationID: req.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", ctx.BaseURL(), shortID),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(responses)
}

// HandleAPIDeleteBatch Delete multiple URLs in batch
// @Summary Delete multiple URLs in a single request
// @Description Accepts a batch of URLs
// @Tags API
// @Accept json
// @Param request body dto.DeleteURLsRequestDTO true "Array of URLs to delete"
// @Success 202 {array}
// @Failure 400 {object} map[string]string "When request body is invalid or empty"
// @Failure 500 {object} map[string]string "When internal server error occurs"
// @Router /api/user/urls [post]
func (c *FiberURLController) HandleAPIDeleteBatch(ctx *fiber.Ctx) error {
	var batchRequestDTO dto.DeleteURLsRequestDTO
	if err := ctx.BodyParser(&batchRequestDTO.URLIDs); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": constants.MsgFailedToParseBody,
		})
	}

	if len(batchRequestDTO.URLIDs) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty batch",
		})
	}

	var err error
	if len(batchRequestDTO.URLIDs) > 1 {
		err = c.service.ConcurrentBatchDelete(ctx.Context(), batchRequestDTO.URLIDs)
	} else {
		err = c.service.BatchDeleteURLs(ctx.Context(), batchRequestDTO.URLIDs)
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to process batch delete",
		})
	}
	return ctx.SendStatus(fiber.StatusAccepted)
}
