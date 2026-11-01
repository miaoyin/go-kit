package lifecycle

import (
	"context"
	"errors"
)

var ErrExternalAlreadyRunning = errors.New("external adapter: already running, start not allowed")

// Adapter 任意对象
type Adapter struct {
	*LifeCycler
	target any
}

func NewAdapter(target any, hook Hook) *Adapter {
	return &Adapter{
		LifeCycler: NewLifeCycler(hook),
		target:     target,
	}
}

// ExternalLifecycleAdapter 只有关闭的外部组件
type ExternalLifecycleAdapter struct {
	inner *LifeCycler
}

func NewExternalLifecycleAdapter(onClose func(ctx context.Context) error, opts ...Option) *ExternalLifecycleAdapter {
	// 内部复用 LifeCycler
	inner := NewLifeCycler(Hook{
		OnStart: func(ctx context.Context) error {
			return ErrExternalAlreadyRunning
		},
		OnClose: onClose,
	},
		append([]Option{
			WithInitState(StateRunning),
			WithRestartable(false),
		}, opts...)...,
	)
	return &ExternalLifecycleAdapter{inner: inner}
}

func (e *ExternalLifecycleAdapter) Start(ctx context.Context) error {
	return e.inner.Start(ctx)
}

func (e *ExternalLifecycleAdapter) Close(ctx context.Context) error {
	return e.inner.Close(ctx)
}

func (e *ExternalLifecycleAdapter) State() State {
	return e.inner.State()
}

func (e *ExternalLifecycleAdapter) IsRunning() bool {
	return e.inner.IsRunning()
}

var _ Lifecycle = (*ExternalLifecycleAdapter)(nil)

func IsExternal(l Lifecycle) bool {
	_, ok := l.(*ExternalLifecycleAdapter)
	return ok
}
