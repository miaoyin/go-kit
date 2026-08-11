package poolutil

import (
	"sync/atomic"
)

//OwnedBuffer 带有引用计数的buffer, 引用数量等于0时释放
type OwnedBuffer struct{
	//内容
	data []byte
	//计数
	ref  int32
	//缓存
	pool *BufferPool
}

func NewOwnedBuffer(pool *BufferPool, src []byte) *OwnedBuffer{
	dst := pool.Get()
	dst = append(dst[:0], src...)
	return &OwnedBuffer{
		data: dst,
		ref: 1,
		pool: pool,
	}
}

//Retain 计数加1
func (b *OwnedBuffer) Retain() {
	if b == nil {
		return
	}
	atomic.AddInt32(&b.ref, 1)
}

//Release 释放
func (b *OwnedBuffer) Release() {
	if b == nil {
		return
	}
	newRef := atomic.AddInt32(&b.ref, -1)
	if newRef == 0 {
		//释放引用
		if b.pool!=nil && b.data != nil {
			b.pool.Put(b.data)
		}
		b.data = nil
		b.pool = nil
	}
}

func (b *OwnedBuffer) Data() []byte {
	if b == nil {
		return nil
	}
	return b.data
}

func (b *OwnedBuffer) Clone() *OwnedBuffer {
	if b == nil {
		return nil
	}
	b.Retain()
	return b
}
