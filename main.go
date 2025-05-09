package main

import (
    "pubsubgrpc/internal/models"

    "pubsubgrpc/internal/config"
    "pubsubgrpc/internal/logger"
    "pubsubgrpc/internal/grpcserver"
)

func main() {
    cfg := config.Load()
    log := logger.New()
    serverSettings := models.ServerSettings{
        Cfg: cfg,
        Log: log,
    }

    err := grpcserver.StartGRPCServer(serverSettings)
    if err != nil{
        log.Fatalf("server error: %v\n", err)
    }
}
