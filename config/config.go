package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
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
}

func (c *Config) setDefaults() {
	if c.ServerAddress == "" {
		c.ServerAddress = ":8080"
	}
	if c.BaseURL == "" {
		c.BaseURL = "http://localhost"
	}
}
