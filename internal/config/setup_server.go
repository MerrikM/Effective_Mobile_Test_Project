package config

import (
	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

type AppConfig struct {
	DatabaseConfig DatabaseConfig `yaml:"databaseConfig"`
	ServerAddr     string         `yaml:"serverAddr"`
}

func LoadConfig(path string) (*AppConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SetupDatabase(dsn string) (*Database, error) {
	return NewDatabaseConnection("postgres", dsn)
}

func SetupRestServer(addr string) (*http.Server, *chi.Mux) {
	router := chi.NewRouter()

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return server, router
}
