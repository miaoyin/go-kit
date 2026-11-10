package atomicx

import "time"

type SnapshotMeta interface {
	Version() uint64
	ChangedAt() time.Time
}

//------------------------------------
//  基于atomic.Pointer类型
//------------------------------------

// PointerSnapshotter 只读快照接口
type PointerSnapshotter[T any] interface {
	SnapshotMeta
	GetSnapshot() (*T, error)
}

// PointerUpdater 替换接口
type PointerUpdater[T any] interface {
	Swap(newVal *T) (oldVal *T)
}

//------------------------------------
//  基于atomic.Value类型
//------------------------------------

//ValueSnapshotter 快照
type ValueSnapshotter[T any] interface {
	SnapshotMeta
	GetSnapshot() (T, error)
}

// ValueUpdater 替换接口
type ValueUpdater[T any] interface {
	Swap(newVal T) (oldVal T)
}

