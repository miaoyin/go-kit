package zapx

import (
    "io"
    "sync/atomic"

    "gopkg.in/natefinch/lumberjack.v2"
)

// SafeLumberjack 包装 lumberjack，增加引用计数和关闭标记
type SafeLumberjack struct {
    lj       *lumberjack.Logger
    ref      int32
    closed   int32
}

func NewSafeLumberjack(lj *lumberjack.Logger) *SafeLumberjack {
    return &SafeLumberjack{
        lj:  lj,
        ref: 1,
    }
}

func (s *SafeLumberjack) AddRef() {
    atomic.AddInt32(&s.ref, 1)
}

func (s *SafeLumberjack) DecRef() {
    if atomic.AddInt32(&s.ref, -1) == 0 {
        _ = s.lj.Close()
        atomic.StoreInt32(&s.closed, 1)
    }
}

func (s *SafeLumberjack) Write(p []byte) (n int, err error) {
    if atomic.LoadInt32(&s.closed) == 1 {
        return 0, io.EOF
    }
    return s.lj.Write(p)
}

// Close 外部主动关闭
func (s *SafeLumberjack) Close() error {
    if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
        return nil
    }
    return s.lj.Close()
}
