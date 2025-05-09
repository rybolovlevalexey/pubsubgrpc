package config

import (
	"os"
	"pubsubgrpc/internal/models"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() *models.Config{
	godotenv.Load()
	grpcPort, _ := strconv.Atoi(os.Getenv("PubSubgRPCPort"))

	config := &models.Config{
		PubSubgRPCPort: grpcPort,
	}

	return config
}