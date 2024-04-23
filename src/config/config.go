package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host           string
	Port           string
	AccessTokenKey string
	APIKey         string
	SessionKey     string
	BackendURL     string
}

func LoadEnv() *Config {
	if len(os.Args) < 2 {
		log.Println("load default env")
		godotenv.Load()
	} else {
		file := os.Args[1]
		log.Println("load", file)
		godotenv.Load(file)
	}

	cfg := &Config{}

	cfg.Host = os.Getenv("HOST")
	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}

	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "3000"
	}

	cfg.AccessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
	if cfg.AccessTokenKey == "" {
		log.Fatal("access token key is required")
	}

	cfg.SessionKey = os.Getenv("SESSION_KEY")
	if cfg.SessionKey == "" {
		log.Fatal("session key is required")
	}

	cfg.BackendURL = os.Getenv("BACKEND_URL")
	if cfg.BackendURL == "" {
		log.Fatal("backend url is required")
	}

	cfg.APIKey = os.Getenv("API_KEY")

	return cfg
}
