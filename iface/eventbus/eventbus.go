package eventbus

import "sync"

//-----------------------------
//  订阅接口
//-----------------------------

type Consumer[T any] interface {
	//Deliver 投递事件, 不能阻塞, 需要处理异步队列
	Deliver(evt *T)
}

type EventBus[T any] struct{
	consumers []Consumer[T]
	mu        sync.RWMutex
}

func (b *EventBus[T]) Subscribe(c Consumer[T]) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.consumers = append(b.consumers, c)
}

func (b *EventBus[T]) Publish(evt *T) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, c := range b.consumers {
		c.Deliver(evt)
	}
}