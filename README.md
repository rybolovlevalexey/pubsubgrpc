# Publisher-Subscriber gRPC server

## Prerequisites
- Go (1.20+)

## Launch project
1. ```git clone https://github.com/rybolovlevalexey/pubsubgrpc```
2. ```go mod download```
3. ```go run main.go```


## Packages description
- subpub - allows to: subscribe to events, receive events when publishing, do not depend on slow subscribers, guarantee the message order, shut down correctly
- proto - grpc proto interface and generated interfaces
- core - config, logger, etc
- models - main structs of the project
- grpcserver - gRPC server which implements publisher-subscriber and tests for this functionality
