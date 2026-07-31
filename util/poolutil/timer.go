package poolutil

import (
	"sync"
	"time"
)

//TimerPool timer对象池
type TimerPool struct {
	pool sync.Pool
}

func NewTimerPool()*TimerPool{
	return &TimerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return time.NewTimer(0)
			},
		},
	}
}

func (p *TimerPool) Get(d time.Duration) *time.Timer {
	t := p.pool.Get().(*time.Timer)
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(d)
	return t
}

func (p *TimerPool) Put(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	p.pool.Put(t)
}
