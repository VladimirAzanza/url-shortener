package server

import (
	"context"

	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func NewFiberServer(urlController *controller.FiberURLController) *fiber.App {
	app := fiber.New()

	app.Get("/:id", urlController.HandleGet)
	app.Post("/", urlController.HandlePost)
	return app
}

func StartFiberServer(lc fx.Lifecycle, app *fiber.App) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(":8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
