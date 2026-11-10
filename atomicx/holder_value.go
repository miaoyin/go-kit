package atomicx

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	_ ValueSnapshotter[any] = (*ValueHolder[any])(nil)
	_ ValueUpdater[any] = (*ValueHolder[any])(nil)
)

//ValueHolder 适合interface类型
type ValueHolder[T any] struct {
	val atomic.Value
	ver       uint64
	changedAt atomic.Pointer[time.Time]
}

func NewValueHolder[T any](initVal ...T) *ValueHolder[T]{
	h := &ValueHolder[T]{}
	if len(initVal) >0 {
		h.Swap(initVal[0])
	}
	return h
}

func (h *ValueHolder[T]) Get() T{
	v, _ := h.val.Load().(T)
	return v
}

func (h *ValueHolder[T]) GetSnapshot() (T, error) {
	v, ok := h.val.Load().(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("value is nil")
	}
	return v, nil
}

func (h *ValueHolder[T]) Version() uint64 {
	return atomic.LoadUint64(&h.ver)
}

func (h *ValueHolder[T]) ChangedAt() time.Time {
	p := h.changedAt.Load()
	if p == nil {
		return time.Time{}
	}
	return *p
}

func (h *ValueHolder[T]) Swap(newVal T) T {
	old := h.val.Swap(newVal)
	atomic.AddUint64(&h.ver, 1)
	now := time.Now()
	h.changedAt.Store(&now)
	oldVal, _ := old.(T)
	return oldVal
}
