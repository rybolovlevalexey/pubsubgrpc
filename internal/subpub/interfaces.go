package subpub

import(
	"context"
)

// обработчик события, передаётся при подписке
type MessageHandler func(msg interface{})

// подписчик имеет возможность отписаться
type Subscription interface {
	Unsubscribe()
}

// 
type SubPub interface {
	Subscribe(subject string, cb MessageHandler) (Subscription, error)  // подписка по ключу
	Publish(subject string, msg interface{}) error  // отправка сообщения
	Close(ctx context.Context) error  // завершение работы
}