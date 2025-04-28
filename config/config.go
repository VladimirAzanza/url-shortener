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
}
