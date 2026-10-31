package lifecycle

import (
    "errors"
    "sync"
)

//ResourceScope 资源作用域
// 管理临时资源
type ResourceScope struct{
    mu      sync.Mutex
    closers []func() error
}

func NewResourceScope() *ResourceScope {
    return &ResourceScope{}
}

//Register 注册资源, 失败回滚
func (s *ResourceScope) Register(fn func() (closeFn func() error, err error)) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    closeFn, err := fn()
    if err != nil {
        s.rollbackLocked()
        return err
    }
    s.closers = append(s.closers, closeFn)
    return nil
}

//rollbackLocked 回滚
func (s *ResourceScope) rollbackLocked() {
    for i := len(s.closers) - 1; i >= 0; i-- {
        _ = s.closers[i]()
    }
    s.closers = nil
}

func (s *ResourceScope) Close() error {
    s.mu.Lock()
    defer s.mu.Unlock()
    var errs []error
    for i := len(s.closers) - 1; i >= 0; i-- {
        if err := s.closers[i](); err != nil {
            errs = append(errs, err)
        }
    }
    s.closers = nil
    if len(errs) > 0 {
        return errors.Join(errs...)
    }
    return nil
}
