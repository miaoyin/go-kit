package poolutil

import "sync"

//RespPool 创建chan用于等待返回结果
// 例如: 消息中带上channel,用于接收返回结果
type RespPool[T any] struct{
	pool sync.Pool
}

func NewRespPool[T any]() *RespPool[T]{
	return &RespPool[T]{pool: sync.Pool{
		New: func() interface{} {
			return make(chan T, 1)
		},
	}}
}

func (p *RespPool[T]) Get() chan T{
	return p.pool.Get().(chan T)
}

//Put 回收
func (p *RespPool[T]) Put(ch chan T) {
	//排空数据
	select {
	case <-ch:
	default:
	}
	p.pool.Put(ch)
}