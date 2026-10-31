package eventbus

import (
    "sync"
    "sync/atomic"
)


type CallbackBus[T any] struct {
    mu        sync.RWMutex
    nextID     uint64
    subscribers map[uint64]func(*T)
}

func NewCallbackBus[T any]()*CallbackBus[T]{
    return &CallbackBus[T]{
        subscribers: make(map[uint64]func(*T)),
    }
}

func (b *CallbackBus[T]) Subscribe(fn func(*T)) uint64{
    if fn == nil {
        return 0
    }
    id := atomic.AddUint64(&b.nextID, 1)
    b.mu.Lock()
    defer b.mu.Unlock()
    b.subscribers[id] = fn
    return id
}

func (b *CallbackBus[T]) Unsubscribe(id uint64) {
    if id == 0 {
        return
    }
    b.mu.Lock()
    defer b.mu.Unlock()
    delete(b.subscribers, id)
}

func (b *CallbackBus[T]) Publish(evt *T) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    for _, fn := range b.subscribers {
        go fn(evt)
    }
}