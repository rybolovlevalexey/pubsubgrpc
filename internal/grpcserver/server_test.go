package grpcserver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	pb "pubsubgrpc/internal/proto"
	"google.golang.org/grpc/metadata"
	"pubsubgrpc/internal/config"
	"pubsubgrpc/internal/logger"
	"pubsubgrpc/internal/models"
)

// заглушки
type mockStream struct {
	recv chan *pb.Event
	ctx  context.Context
}

func (m *mockStream) Send(e *pb.Event) error {
	m.recv <- e
	return nil
}

func (m *mockStream) Context() context.Context {
	return m.ctx
}

func (m *mockStream) SetHeader(md metadata.MD) error  { return nil }
func (m *mockStream) SendHeader(md metadata.MD) error { return nil }
func (m *mockStream) SetTrailer(md metadata.MD)       {}
func (m *mockStream) SendMsg(interface{}) error       { return nil }
func (m *mockStream) RecvMsg(interface{}) error       { return nil }


// тесты
func TestPublishSubscribe(t *testing.T) {
	cfg := config.Load()
	log := logger.New()
	serverSettings := models.ServerSettings{Cfg: cfg, Log: log}
	srv := NewPubSubServer(serverSettings).(*pubSubServer)

	// Контекст не отменяется
	ctx := context.Background()
	msgChan := make(chan *pb.Event, 1)
	stream := &mockStream{recv: msgChan, ctx: ctx}

	go func() {
		err := srv.Subscribe(&pb.SubscribeRequest{Key: "test"}, stream)
		assert.NoError(t, err)
	}()

	time.Sleep(50 * time.Millisecond)

	_, err := srv.Publish(context.Background(), &pb.PublishRequest{Key: "test", Data: "hello"})
	assert.NoError(t, err)

	select {
	case msg := <-msgChan:
		assert.Equal(t, "hello", msg.Data)
	case <-time.After(1 * time.Second):
		t.Fatal("no message received")
	}
}

func TestSubscribeAndPublishWithContextCancel(t *testing.T) {
	cfg := config.Load()
	log := logger.New()
	serverSettings := models.ServerSettings{Cfg: cfg, Log: log}
	srv := NewPubSubServer(serverSettings).(*pubSubServer)

	msgChan := make(chan *pb.Event, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	stream := &mockStream{recv: msgChan, ctx: ctx}

	done := make(chan error)

	go func() {
		err := srv.Subscribe(&pb.SubscribeRequest{Key: "weather"}, stream)
		done <- err
	}()

	// Ждём, чтобы подписка установилась
	time.Sleep(20 * time.Millisecond)

	// Публикуем событие
	_, err := srv.Publish(context.Background(), &pb.PublishRequest{
		Key:  "weather",
		Data: "sunny",
	})
	assert.NoError(t, err)

	select {
	case msg := <-msgChan:
		assert.Equal(t, "sunny", msg.Data)
	case <-time.After(1 * time.Second):
		t.Fatal("message not received")
	}

	// Ждём завершения подписки по таймауту
	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("subscribe did not exit on context cancel")
	}
}
