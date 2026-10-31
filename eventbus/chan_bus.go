package eventbus

import (
    "sync"
    "sync/atomic"
)

type ChanBus[T any] struct {
    mu          sync.RWMutex
    subscribers map[chan *T]struct{}
    //统计丢失数
    dropCounter uint64
    //onDrop 队列满了回调, 可打印日志
    onDrop func(evt *T)
}

func NewChanBus[T any](onDrop func(*T)) *ChanBus[T] {
    return &ChanBus[T]{
        subscribers: make(map[chan *T]struct{}),
        onDrop: onDrop,
    }
}

func (b *ChanBus[T]) Subscribe(bufSize int) chan *T {
    //防止慢订阅
    ch := make(chan *T, bufSize)
    b.mu.Lock()
    defer b.mu.Unlock()
    b.subscribers[ch] = struct{}{}
    return ch
}

// Unsubscribe 取消订阅
func (b *ChanBus[T]) Unsubscribe(ch chan *T) {
    b.mu.Lock()
    defer b.mu.Unlock()
    delete(b.subscribers, ch)
    close(ch)
}

func (b *ChanBus[T]) DropCount() uint64 {
    return atomic.LoadUint64(&b.dropCounter)
}

func (b *ChanBus[T]) Publish(e *T) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    for ch := range b.subscribers {
        select {
        case ch <- e:
        default:
            // 消息丢弃
            atomic.AddUint64(&b.dropCounter, 1)
            if b.onDrop != nil {
                go b.onDrop(e)
            }
        }
    }
}
