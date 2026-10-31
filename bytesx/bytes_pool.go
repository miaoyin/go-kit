package bytesx

import "sync"

const maxByteSliceCap = 64 * 1024

// BytesPool 字节缓存
type BytesPool struct {
	pool *sync.Pool
}

func NewBytesPool(size int) *BytesPool {
	return &BytesPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, size)
			},
		},
	}
}

func (p *BytesPool) Get() []byte {
	return p.pool.Get().([]byte)[:0]
}

func (p *BytesPool) Put(buf []byte) {
	if cap(buf) > maxByteSliceCap {
		return
	}
	p.pool.Put(buf[:0])
}

func (p *BytesPool) NewOwnedBuffer(src []byte) *OwnedBuffer {
	dst := p.pool.Get().([]byte)
	dst = append(dst[:0], src...)
	return &OwnedBuffer{
		data: dst,
		ref:  1,
		pool: p,
	}
}
