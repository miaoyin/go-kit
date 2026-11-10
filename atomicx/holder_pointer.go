package atomicx

import (
	"errors"
	"sync/atomic"
	"time"
)


var (
	_ PointerSnapshotter[any] = (*PointerHolder[any])(nil)
	_ PointerUpdater[any] = (*PointerHolder[any])(nil)
)

//PointerHolder 指针持有者
// 适合struct的泛型
type PointerHolder[T any] struct {
	ptr atomic.Pointer[T]
	//不是业务一致性, 也不要加sync.RWMutex, 会破坏Pointer优势
	ver       uint64
	changedAt atomic.Pointer[time.Time]
}

func NewPointerHolder[T any](initVal ...*T) *PointerHolder[T] {
	h := &PointerHolder[T]{}
	if len(initVal) > 0 {
		h.Swap(initVal[0])
	}
	return h
}

func (h *PointerHolder[T]) Get() *T{
	return h.ptr.Load()
}

func (h *PointerHolder[T]) GetSnapshot() (*T, error) {
	p := h.ptr.Load()
	if p == nil {
		return nil, errors.New("holder: value not ready")
	}
	return p, nil
}

func (h *PointerHolder[T]) Swap(newVal *T) (old *T) {
	oldPtr := h.ptr.Swap(newVal)
	atomic.AddUint64(&h.ver, 1)
	now := time.Now()
	h.changedAt.Store(&now)
	return oldPtr
}

func (h *PointerHolder[T]) Version() uint64 {
	return atomic.LoadUint64(&h.ver)
}

func (h *PointerHolder[T]) ChangedAt() time.Time {
	p := h.changedAt.Load()
	if p == nil {
		return time.Time{}
	}
	return *p
}
