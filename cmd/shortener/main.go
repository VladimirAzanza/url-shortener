package main

import (
	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"github.com/VladimirAzanza/url-shortener/internal/server"
	"github.com/VladimirAzanza/url-shortener/internal/services"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			services.NewURLService,
			// controller.NewURLController,
			controller.NewFiberURLController,
			//server.NewServer,
			server.NewFiberServer,
		),
		fx.Invoke(server.StartFiberServer),
	).Run()
}
