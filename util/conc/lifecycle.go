package conc

import (
    "context"
    "sync"
)

//Lifecycle 异步任务生命周期
type Lifecycle struct{
    stopOnce sync.Once
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

func NewLifecycle(parentCtx context.Context) *Lifecycle {
    ctx, cancel := context.WithCancel(parentCtx)
    return &Lifecycle{
        ctx:    ctx,
        cancel: cancel,
    }
}

//Spawn 启动子任务
func (l *Lifecycle) Spawn(fn func(ctx context.Context)) {
    l.wg.Add(1)
    go func() {
        defer l.wg.Done()
        fn(l.ctx)
    }()
}

func (l *Lifecycle) Stop() {
    l.stopOnce.Do(func() {
        l.cancel()
        l.wg.Wait()
    })
}

func (l *Lifecycle) Wait() {
    l.wg.Wait()
}

func (l *Lifecycle) Context() context.Context {
    return l.ctx
}
