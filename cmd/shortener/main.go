package main

import (
	"github.com/VladimirAzanza/url-shortener/config"
	_ "github.com/VladimirAzanza/url-shortener/docs"
	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/VladimirAzanza/url-shortener/internal/server"
	"github.com/VladimirAzanza/url-shortener/internal/services"
	_ "github.com/mattn/go-sqlite3"
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
		repo.NewDB,
		services.NewURLService,
		controller.NewFiberURLController,
		server.NewFiberServer,
	),
	fx.Invoke(server.StartFiberServer),
)

// Agregar mock uber y testear los controllers
