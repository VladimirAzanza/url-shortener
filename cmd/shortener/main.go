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
			controller.NewURLController,
			server.NewServer,
		),
		fx.Invoke(server.StartServer),
	).Run()
}
