package main

import (
	"github.com/VladimirAzanza/url-shortener/config"
	_ "github.com/VladimirAzanza/url-shortener/docs"
	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"github.com/VladimirAzanza/url-shortener/internal/server"
	"github.com/VladimirAzanza/url-shortener/internal/services"
	"go.uber.org/fx"
)

// @title URL shortener API
// @version 1.0
// @description This is a sample swagger for URL Shortener
// @contact.email vladimirazanza@gmail.com
// @host localhost:8080
// @BasePath /
func main() {
	fx.New(Module).Run()
}

var Module = fx.Module(
	"main",
	fx.Supply(config.NewConfig()),
	fx.Provide(
		services.NewURLService,
		controller.NewFiberURLController,
		server.NewFiberServer,
	),
	fx.Invoke(server.StartFiberServer),
)
