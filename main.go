package main

import (
    "pubsubgrpc/internal/models"

    "pubsubgrpc/internal/core"
    "pubsubgrpc/internal/grpcserver"
)

func main() {
    cfg := core.LoadConfig()
    log := core.NewLogger()
    serverSettings := models.ServerSettings{
        Cfg: cfg,
        Log: log,
    }

    err := grpcserver.StartGRPCServer(serverSettings)
    if err != nil{
        log.Fatalf("server error: %v\n", err)
    }
}
