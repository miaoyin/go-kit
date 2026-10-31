package atomicx

import (
	"errors"
	"sync/atomic"
	"time"
)

// Snapshotter 只读快照接口
type Snapshotter[T any] interface {
	GetSnapshot() (*T, error)
	Version() uint64
	ChangedAt() time.Time
}

// Updater 替换接口
type Updater[T any] interface {
	Swap(newVal *T) (old *T)
}

type Holder[T any] struct {
	ptr atomic.Pointer[T]
	//不是业务一致性, 也不要加sync.RWMutex, 会破坏Pointer优势
	ver       uint64
	changedAt atomic.Pointer[time.Time]
}

func NewHolder[T any](initVal ...*T) *Holder[T] {
	h := &Holder[T]{}
	if len(initVal) > 0 {
		h.Swap(initVal[0])
	}
	return h
}

func NewHolderSnapshotterAndUpdater[T any](initVal ...*T) (Snapshotter[T], Updater[T]) {
	h := NewHolder(initVal...)
	return h, nil
}

func (h *Holder[T]) GetSnapshot() (*T, error) {
	p := h.ptr.Load()
	if p == nil {
		return nil, errors.New("holder: value not ready")
	}
	return p, nil
}

func (h *Holder[T]) Swap(newVal *T) (old *T) {
	oldPtr := h.ptr.Swap(newVal)
	atomic.AddUint64(&h.ver, 1)
	now := time.Now()
	h.changedAt.Store(&now)
	return oldPtr
}

func (h *Holder[T]) Version() uint64 {
	return atomic.LoadUint64(&h.ver)
}

func (h *Holder[T]) ChangedAt() time.Time {
	p := h.changedAt.Load()
	if p == nil {
		return time.Time{}
	}
	return *p
}
