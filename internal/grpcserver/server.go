package grpcserver


import (
	"context"
	"fmt"
	"net"
	"log"
	"sync"

	"pubsubgrpc/internal/models"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "pubsubgrpc/internal/proto"
)


// StartGRPCServer запускает gRPC сервер
func StartGRPCServer(serverSettings models.ServerSettings) error {
	lst, err := net.Listen("tcp", fmt.Sprintf(":%d", serverSettings.Cfg.PubSubgRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterPubSubServer(s, NewPubSubServer(serverSettings))

	log.Printf("gRPC server listening on :%d\n", serverSettings.Cfg.PubSubgRPCPort)
	return s.Serve(lst)
}


type pubSubServer struct {
	pb.UnimplementedPubSubServer
	models.ServerSettings
	subscribers map[string][]chan string
	mu          sync.Mutex
}

func NewPubSubServer(serverSettings models.ServerSettings) pb.PubSubServer {
	serverSettings.Log.Printf("get request to create new PubSub server\n")

	return &pubSubServer{
		subscribers: make(map[string][]chan string),
		ServerSettings: serverSettings,
	}
}

func (s *pubSubServer) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	s.Log.Printf("pubSubServer subscribe method with param %s\n", req.GetKey())
	key := req.GetKey()
	ch := make(chan string, 10)

	s.mu.Lock()
	s.subscribers[key] = append(s.subscribers[key], ch)
	s.mu.Unlock()
	s.Log.Printf("added new subscriber %s\n", req.GetKey())

	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			// Удаляем подписчика и закрываем канал
			s.mu.Lock()
			subs := s.subscribers[key]
			for i, c := range subs {
				if c == ch {
					s.subscribers[key] = append(subs[:i], subs[i+1:]...)
					break
				}
			}
			s.mu.Unlock()
			return ctx.Err()

		case msg := <-ch:
			if err := stream.Send(&pb.Event{Data: msg}); err != nil {
				return err
			}
		}
	}
}

func (s *pubSubServer) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	s.Log.Printf("pubSubServer Publish method with params %s %s\n", req.GetKey(), req.GetData())
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ch := range s.subscribers[req.GetKey()] {
		select {
		case ch <- req.GetData():
		default:
		}
	}

	return &emptypb.Empty{}, nil
}
