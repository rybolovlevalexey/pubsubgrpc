package config

import (
	"pubsubgrpc/internal/models"
	// "os"

	"github.com/joho/godotenv"
)

func Load() *models.Config{
	godotenv.Load()

	config := &models.Config{  // os.Getenv

	}

	return config
}