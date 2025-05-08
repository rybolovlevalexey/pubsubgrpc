package subpub

import(
	"testing"
	"time"
	"sync"
	"context"
)

func TestSubscribeAndPublish(t *testing.T) {
	sp := NewSubPub()

	var mu sync.Mutex
	received := []interface{}{}

	sub, err := sp.Subscribe("test", func(msg interface{}) {
		mu.Lock()
		defer mu.Unlock()
		received = append(received, msg)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = sp.Publish("test", "message1")
	if err != nil {
		t.Fatal(err)
	}

	err = sp.Publish("test", "message2")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	sub.Unsubscribe()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	sp.Close(ctx)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Errorf("expected 2 messages, got %d", len(received))
	}
	if received[0] != "message1" || received[1] != "message2" {
		t.Errorf("unexpected message order: %v", received)
	}
}

func TestUnsubscribe(t *testing.T) {
	sp := NewSubPub()

	var count int
	sub, _ := sp.Subscribe("a", func(msg interface{}) {
		count++
	})

	sub.Unsubscribe()
	sp.Publish("a", "should be ignored")

	time.Sleep(50 * time.Millisecond)

	if count != 0 {
		t.Errorf("expected 0 messages after unsubscribe, got %d", count)
	}
}

func TestSlowSubscriberDoesNotBlockOthers(t *testing.T) {
	sp := NewSubPub()

	var fastReceived int
	_, _ = sp.Subscribe("s", func(msg interface{}) {
		fastReceived++
	})

	_, _ = sp.Subscribe("s", func(msg interface{}) {
		time.Sleep(1 * time.Second) // slow subscriber
	})

	_ = sp.Publish("s", "data1")
	_ = sp.Publish("s", "data2")

	time.Sleep(200 * time.Millisecond)

	if fastReceived != 2 {
		t.Errorf("fast subscriber didn't receive all messages (got %d)", fastReceived)
	}
}

func TestCloseWithContext(t *testing.T) {
	sp := NewSubPub()

	_, _ = sp.Subscribe("ctx", func(msg interface{}) {
		time.Sleep(time.Second)
	})

	_ = sp.Publish("ctx", "delayed")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := sp.Close(ctx)
	if err == nil {
		t.Errorf("expected context deadline error, got nil")
	}
}