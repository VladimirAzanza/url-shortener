package server

import (
	"context"

	"github.com/VladimirAzanza/url-shortener/config"
	_ "github.com/VladimirAzanza/url-shortener/docs"
	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"github.com/VladimirAzanza/url-shortener/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/swagger"
	"go.uber.org/fx"
)

// ./shortener -a :8081
// SERVER_ADDRESS=:8082 ./shortener
func NewFiberServer(urlController *controller.FiberURLController) *fiber.App {
	app := fiber.New()

	compressionMiddleware := compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})

	app.Use(middleware.MiddlewareZerolog())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/:id", urlController.HandleGet)
	app.Post("/", compressionMiddleware, urlController.HandlePost)

	api := app.Group("/api")
	{
		api.Post("/shorten", compressionMiddleware, urlController.HandleAPIPost)
	}

	return app
}

func StartFiberServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(cfg.ServerAddress)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
