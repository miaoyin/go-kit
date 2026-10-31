package chanx

import "sync"

// SignalQueue 信号队列
//
//	解决收到chan通知时,已经将数据取出
type SignalQueue[T any] struct {
	mu     sync.Mutex
	buf    []T
	notify chan struct{}
	closed bool //避免notify重复close
}

func NewSignalQueue[T any](capHint int) *SignalQueue[T] {
	return &SignalQueue[T]{
		buf:    make([]T, 0, capHint),
		notify: make(chan struct{}, 1),
	}
}

// Enqueue 入队
func (q *SignalQueue[T]) Enqueue(item T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	q.buf = append(q.buf, item)
	//通知
	select {
	case q.notify <- struct{}{}:
	default:
	}
	return true
}

// TryDequeue 出队
func (q *SignalQueue[T]) TryDequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.buf) == 0 {
		var zero T
		return zero, false
	}
	item := q.buf[0]
	q.buf = q.buf[1:]
	return item, true
}

// Peek 看第一个数据
func (q *SignalQueue[T]) Peek() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.buf) == 0 {
		var zero T
		return zero, false
	}
	return q.buf[0], true
}

// Close 关闭
func (q *SignalQueue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	close(q.notify)
}

// Notify 通知
func (q *SignalQueue[T]) Notify() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	//通知
	select {
	case q.notify <- struct{}{}:
	default:
	}
}

// Wait 信号
func (q *SignalQueue[T]) Wait() <-chan struct{} {
	return q.notify
}

func (q *SignalQueue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.buf)
}

// Clear 清空
func (q *SignalQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.buf = q.buf[:0]
	//排空
	select {
	case <-q.notify:
	default:
	}
}
