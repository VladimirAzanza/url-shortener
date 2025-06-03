package main

import (
	"github.com/VladimirAzanza/url-shortener/config"
	_ "github.com/VladimirAzanza/url-shortener/docs"
	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"github.com/VladimirAzanza/url-shortener/internal/repo"
	"github.com/VladimirAzanza/url-shortener/internal/server"
	"github.com/VladimirAzanza/url-shortener/internal/services"
	_ "github.com/lib/pq"
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
// Agregar el guardado de datos en la BD. Agregar un flag para que el usuario elija donde guardar los datos, bD o file.json
// Agregar tests de benchmarking
//  POST /api/shorten/batch -> request (agregar tx.begin, commit, rollback):
// [
//     {
//         "correlation_id": "<строковый идентификатор>",
//         "original_url": "<URL для сокращения>"
//     },
//     ...
// ]

// response:
// 	[
//     {
//         "correlation_id": "<строковый идентификатор из объекта запроса>",
//         "short_url": "<результирующий сокращённый URL>"
//     },
//     ...
// ]

// Intentar usar errors is errores as y join => revisar increment 13
// investigar sobre sqlc
