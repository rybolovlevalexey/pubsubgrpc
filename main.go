package main

import (
    "pubsubgrpc/internal/config"
    "pubsubgrpc/internal/logger"
    "pubsubgrpc/internal/grpcserver"
)

func main() {
    cfg := config.Load()
    log := logger.New()

    grpcserver.Run(cfg, log)
}
