package poolutil

import "sync"


const maxByteSliceCap = 64*1024

//BufferPool 字节缓存
type BufferPool struct{
	pool *sync.Pool
}

func NewBufferPool(size int)*BufferPool{
	return &BufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, size)
			},
		},
	}
}

func (p *BufferPool) Get() []byte{
	return p.pool.Get().([]byte)[:0]
}

func (p *BufferPool) Put(buf []byte) {
	if cap(buf) > maxByteSliceCap {
		return
	}
	p.pool.Put(buf[:0])
}

func (p *BufferPool) NewOwnedBuffer(src []byte) *OwnedBuffer {
	dst := p.pool.Get().([]byte)
	dst = append(dst[:0], src...)
	return &OwnedBuffer{
		data: dst,
		ref: 1,
		pool: p,
	}
}