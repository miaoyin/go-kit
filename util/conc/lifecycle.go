package conc

import (
    "context"
    "sync"
)

//Lifecycle 异步任务生命周期
type Lifecycle struct{
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup
}

func NewLifecycle() *Lifecycle {
    ctx, cancel := context.WithCancel(context.Background())
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
    l.cancel()
    l.wg.Wait()
}

func (l *Lifecycle) Context() context.Context {
    return l.ctx
}
