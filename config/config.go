package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseType    string `env:"DATABASE_TYPE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func NewConfig() *Config {
	cfg := &Config{}

	cfg.loadFromEnv()
	cfg.parseFlags()
	cfg.setDefaults()

	return cfg
}

func (c *Config) parseFlags() {
	flag.StringVar(
		&c.ServerAddress, "a", c.ServerAddress, "Server address (env: SERVER_ADDRESS)",
	)
	flag.StringVar(
		&c.BaseURL, "b", c.BaseURL, "Server address (env: BASE_URL)",
	)
	flag.StringVar(
		&c.FileStoragePath, "f", c.FileStoragePath, "File storage path (env: FILE_STORAGE_PATH)",
	)
	flag.StringVar(
		&c.DatabaseType, "dt", c.DatabaseType, "Database type (sqlite|postgres) (env: DATABASE_TYPE)",
	)
	flag.StringVar(
		&c.DatabaseDSN, "dsn", c.DatabaseDSN, "Database connection string (env: DATABASE_DSN)",
	)
	if hasFlags() {
		flag.Parse()
	}
}

func hasFlags() bool {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-a") {
			return true
		}
		if strings.HasPrefix(arg, "-b") {
			return true
		}
		if strings.HasPrefix(arg, "-f") {
			return true
		}
		if strings.HasPrefix(arg, "-d") {
			return true
		}
		if strings.HasPrefix(arg, "-dsn") {
			return true
		}
	}
	return false
}

func (c *Config) loadFromEnv() {
	if address, exists := os.LookupEnv("SERVER_ADDRESS"); exists {
		c.ServerAddress = address
	}
	if baseURL, exists := os.LookupEnv("BASE_URL"); exists {
		c.BaseURL = baseURL
	}
	if filePath, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
		c.FileStoragePath = filePath
	}
	if dbType, exists := os.LookupEnv("DATABASE_TYPE"); exists {
		c.DatabaseType = dbType
	}
	if dbDSN, exists := os.LookupEnv("DATABASE_DSN"); exists {
		c.DatabaseDSN = dbDSN
	}
}

func (c *Config) setDefaults() {
	if c.ServerAddress == "" {
		c.ServerAddress = ":8080"
	}
	if c.BaseURL == "" {
		c.BaseURL = "http://localhost"
	}
	if c.FileStoragePath == "" {
		c.FileStoragePath = "/tmp/short-url-db.json"
	}
	if c.DatabaseType == "" {
		c.DatabaseType = "sqlite"
	}
	if c.DatabaseDSN == "" {
		if c.DatabaseType == "sqlite" {
			c.DatabaseDSN = "file:urlshortener.db?cache=shared&mode=rwc"
		} else {
			c.DatabaseDSN = "host=localhost user=postgres password=postgres dbname=urlshortener port=5432 sslmode=disable"
		}
	}
}
