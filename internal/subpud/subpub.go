package subpud

import (
	"context"
	"errors"
)


func NewSubPub() SubPub {
	return &subPubImpl{
		subs:    make(map[string]map[*subscriber]struct{}),
		closeCh: make(chan struct{}),
	}
}

func (s *subPubImpl) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, errors.New("subpub is closed")
	}

	sub := &subscriber{
		cb:     cb,
		ch:     make(chan interface{}, 64),
		closed: make(chan struct{}),
	}

	if s.subs[subject] == nil {
		s.subs[subject] = make(map[*subscriber]struct{})
	}
	s.subs[subject][sub] = struct{}{}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case msg := <-sub.ch:
				sub.cb(msg)
			case <-sub.closed:
				return
			case <-s.closeCh:
				return
			}
		}
	}()

	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.subs[subject], sub)
		close(sub.closed)
	}

	return &subscription{unsubscribe: unsubscribe}, nil
}

func (s *subPubImpl) Publish(subject string, msg interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return errors.New("subpub is closed")
	}

	for sub := range s.subs[subject] {
		select {
		case sub.ch <- msg:
		default: // Не блокируем — медленный подписчик не тормозит остальных
		}
	}
	return nil
}

func (s *subPubImpl) Close(ctx context.Context) error {
	s.mu.Lock()
	s.closed = true
	close(s.closeCh)
	s.mu.Unlock()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
