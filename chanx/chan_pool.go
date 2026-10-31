package chanx

import "sync"

// ChanPool 创建chan用于等待返回结果
// 例如: 消息中带上channel,用于接收返回结果
type ChanPool[T any] struct {
	pool sync.Pool
}

func NewChanPool[T any](size int) *ChanPool[T] {
	return &ChanPool[T]{pool: sync.Pool{
		New: func() interface{} {
			return make(chan T, size)
		},
	}}
}

func (p *ChanPool[T]) Get() chan T {
	return p.pool.Get().(chan T)
}

// Put 回收
func (p *ChanPool[T]) Put(ch chan T) {
	//排空数据
	isEmpty := false
	for isEmpty == false {
		select {
		case <-ch:
			isEmpty = false
		default:
			isEmpty = true
		}
	}
	p.pool.Put(ch)
}
