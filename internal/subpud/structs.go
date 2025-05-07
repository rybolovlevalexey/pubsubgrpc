package subpud

import (
	"sync"
)

type subscriber struct {
	cb     MessageHandler
	ch     chan interface{}
	closed chan struct{}
}

type subscription struct {
	unsubscribe func()
}

func (s *subscription) Unsubscribe() {
	s.unsubscribe()
}

type subPubImpl struct {
	mu         sync.RWMutex
	subs       map[string]map[*subscriber]struct{}
	closed     bool
	closeCh    chan struct{}
	wg         sync.WaitGroup
}