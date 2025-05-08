package subpub

import (
	"sync"
)

// подписчик
type subscriber struct {
	cb     MessageHandler  // функция обработки сообщений
	ch     chan interface{}  // канал для получения сообщений
	closed chan struct{}  // канал для закрытия
}

// источник сообщений
type subscription struct {
	unsubscribe func()
}


type subPubImpl struct {
	mu         sync.RWMutex
	subs       map[string]map[*subscriber]struct{}  // словарь подписчиков по ключам
	closed     bool
	closeCh    chan struct{}
	wg         sync.WaitGroup
}