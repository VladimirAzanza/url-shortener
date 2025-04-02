package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func MiddlewareZerolog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		log.Info().
			Str("method", c.Method()).
			Str("uri", c.Path()).
			Str("duration", duration.String()).
			Int("status", c.Response().StatusCode()).
			Int("response length", len(c.Response().Body())).
			Msg("request processed")

		return err
	}
}
