package server

import (
	"context"
	"time"

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
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		//  Max body limit set to 10 MB
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(middleware.MiddlewareZerolog())

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/ping", urlController.GetDBPing)
	app.Get("/:id", urlController.HandleGet)
	app.Post("/", urlController.HandlePost)

	api := app.Group("/api")
	{
		api.Post("/shorten", urlController.HandleAPIPost)
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
