package lifecycle

import (
	"context"
	"sync"
)

// Spawner 异步任务生命周期
// 不可复用, 在Start时新建
type Spawner struct {
	*LifeCycler
	onceStop sync.Once
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewSpawner(parentCtx context.Context) *Spawner {
	ctx, cancel := context.WithCancel(parentCtx)
	sp := &Spawner{
		ctx:    ctx,
		cancel: cancel,
	}
	sp.LifeCycler = NewLifeCycler(Hook{
		OnClose: sp.onClose,
	})
	_ = sp.Start(ctx)
	return sp
}

// Spawn 启动子任务
func (l *Spawner) Spawn(fn func(ctx context.Context)) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		fn(l.ctx)
	}()
}

func (l *Spawner) onClose(_ context.Context) error {
	l.onceStop.Do(func() {
		l.cancel()
		l.wg.Wait()
	})
	return nil
}

func (l *Spawner) Wait() {
	l.wg.Wait()
}

// Context 上下文
// Lifecycle一次性, 不直接持有context, 使用函数包装
func (l *Spawner) Context() context.Context {
	return l.ctx
}
