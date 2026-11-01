package lifecycle

import (
	"context"
	"errors"
	"sync/atomic"
)

var ErrAlreadyClosed = errors.New("lifecycle: already closed")
var ErrNotRunning = errors.New("lifecycle: not running")

// State 模块状态
type State uint8

const (
	//StateIdle 已创建，未启动
	StateIdle State = iota
	//StateRunning 正常运行, running时持续状态, Started瞬态不要用started
	StateRunning
	//StateClosing 关闭中, 临时状态
	StateClosing
	//StateClosed 已关闭（终态，无法重启）
	StateClosed
)

func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateRunning:
		return "running"
	case StateClosed:
		return "closed"
	default:
		return "unknown"
	}
}

type Hook struct {
	OnStart func(ctx context.Context) error
	OnClose func(ctx context.Context) error
}

type Lifecycle interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	State() State
	IsRunning() bool
}

//-------------------------------------------
//  LifeCycler 具体实现
//-------------------------------------------

type Option func(*LifeCycler)

// WithRestartable 允许重启
func WithRestartable(v bool) Option {
	return func(l *LifeCycler) {
		l.restartable = v
	}
}

// WithStateCallback 状态更新
func WithStateCallback(fn func(old, new State)) Option {
	return func(l *LifeCycler) {
		l.stateChange = fn
	}
}

// WithInitState 初始状态
func WithInitState(s State) Option {
	return func(l *LifeCycler) {
		l.state.Store(uint32(s))
	}
}

type LifeCycler struct {
	state atomic.Uint32
	hook  Hook
	//是否能重启
	restartable bool
	//状态改变, 可以写日志
	stateChange func(old, new State)
}

func NewLifeCycler(hook Hook, opts ...Option) *LifeCycler {
	l := &LifeCycler{
		hook:        hook,
		restartable: false,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *LifeCycler) Start(ctx context.Context) error {
	oldState := StateIdle
	newState := StateRunning
	if !l.state.CompareAndSwap(uint32(oldState), uint32(newState)) {
		st := State(l.state.Load())
		return errors.New("lifecycle: already in state " + st.String())
	}
	//状态回调
	if l.stateChange != nil {
		l.stateChange(oldState, newState)
	}

	//启动
	if l.hook.OnStart != nil {
		if err := l.hook.OnStart(ctx); err != nil {
			//状态回调
			rollbackOld := StateRunning
			rollbackNew := StateIdle
			l.state.CompareAndSwap(uint32(rollbackOld), uint32(rollbackNew))
			if l.stateChange != nil {
				l.stateChange(rollbackOld, rollbackNew)
			}
			return err
		}
	}
	return nil
}

func (l *LifeCycler) Close(ctx context.Context) error {
	oldState := StateRunning
	newState := StateClosing
	if !l.state.CompareAndSwap(uint32(StateRunning), uint32(StateClosing)) {
		switch State(l.state.Load()) {
		case StateClosing, StateClosed:
			return nil
		default:
			return ErrNotRunning
		}
	}
	defer func() {
		var finalState State
		if l.restartable {
			finalState = StateIdle //允许重启
		} else {
			finalState = StateClosed //不允许重启
		}
		l.state.Store(uint32(finalState))
		//状态回调
		if l.stateChange != nil {
			l.stateChange(newState, finalState)
		}
	}()
	//状态回调
	if l.stateChange != nil {
		l.stateChange(oldState, newState)
	}

	//关闭
	if l.hook.OnClose != nil {
		return l.hook.OnClose(ctx)
	}
	return nil
}

func (l *LifeCycler) State() State {
	return State(l.state.Load())
}

func (l *LifeCycler) IsRunning() bool {
	return State(l.state.Load()) == StateRunning
}

// NamedLifeCycler 带名字
type NamedLifeCycler struct {
	*LifeCycler
	name string
}

func NewNamedLifeCycler(name string, hook Hook, opts ...Option) *NamedLifeCycler {
	return &NamedLifeCycler{
		LifeCycler: NewLifeCycler(hook, opts...),
		name:       name,
	}
}

func (l *NamedLifeCycler) Name() string {
	return l.name
}
