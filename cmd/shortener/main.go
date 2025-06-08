package main

import (
	"github.com/VladimirAzanza/url-shortener/config"
	_ "github.com/VladimirAzanza/url-shortener/docs"
	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"github.com/VladimirAzanza/url-shortener/internal/repo"
	filerepo "github.com/VladimirAzanza/url-shortener/internal/repo/file_repo"
	"github.com/VladimirAzanza/url-shortener/internal/repo/memory"
	"github.com/VladimirAzanza/url-shortener/internal/repo/postgres"
	"github.com/VladimirAzanza/url-shortener/internal/repo/sqlite"
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
		provideRepository,
		services.NewURLService,
		controller.NewFiberURLController,
		server.NewFiberServer,
	),
	fx.Invoke(server.StartFiberServer),
)

func provideRepository(cfg *config.Config) repo.IURLRepository {
	switch cfg.StorageType {
	case "memory":
		return memory.NewMemoryRepository()
	case "file":
		return filerepo.NewFileRepository(cfg)
	case "sqlite":
		db, _ := repo.NewDB(cfg)
		return sqlite.NewSQLiteRepository(db)
	case "postgres":
		db, _ := repo.NewDB(cfg)
		return postgres.NewPostgreSQLRepository(db)
	default:
		panic("unsupported storage type")
	}
}

// Agregar mock uber y testear los controllers
// Agregar el guardado de datos en la BD. Agregar un flag para que el usuario elija donde guardar los datos, bD o file.json
// Agregar tests de benchmarking
// Intentar usar errors is errores as y join => revisar increment 13
// investigar sobre sqlc
// agregar autentificacion
// agregar postgres
