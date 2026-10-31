package lifecycle

import (
    "context"
    "sync"
)

//Lifecycle 异步任务生命周期
// 不可复用, 在Start时新建
type Lifecycle struct{
    onceStop sync.Once
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
    l.onceStop.Do(func() {
        l.cancel()
        l.wg.Wait()
    })
}

func (l *Lifecycle) Wait() {
    l.wg.Wait()
}

//Context 上下文
// Lifecycle一次性, 不直接持有context, 使用函数包装
func (l *Lifecycle) Context() context.Context {
    return l.ctx
}
