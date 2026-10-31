package module

import (
	"sync/atomic"
)

// EmptyConfig 无参数配置
// 例如 NewBaseModule[EmptyConfig]("name", EmptyConfig{})
type EmptyConfig struct{}

// BaseModule 基础模块
// 避免重复接口
type BaseModule[T any] struct {
	name  string
	state atomic.Uint32
	//confVal 配置, 支持热更新
	confVal atomic.Pointer[T]
	//confVer 配置版本, 默认0
	confVer atomic.Int64
}

func NewBaseModule[T any](name string, initCfg T) *BaseModule[T] {
	b := &BaseModule[T]{
		name: name,
	}
	b.confVal.Store(&initCfg)
	b.state.Store(uint32(StateInit))
	return b
}

func (b *BaseModule[T]) Name() string {
	return b.name
}

// GetConfig 读取配置, 返回父辈
// 1.高频调用时, 使用缓存, 定时刷新减少调用频率
// 2.configVer版本号, 区分是否更新
func (b *BaseModule[T]) GetConfig() (*T, int64) {
	return b.confVal.Load(), b.confVer.Load()
}

// Config 读取配置
// 读取后不要修改, 修改使用SetConfig
func (b *BaseModule[T]) Config() *T {
	return b.confVal.Load()
}

// SetConfig 设置配置
// 根据需要覆盖
func (b *BaseModule[T]) SetConfig(cfg T) error {
	b.confVal.Store(&cfg)
	b.confVer.Add(1)
	return nil
}

func (b *BaseModule[T]) State() State {
	return State(b.state.Load())
}

func (b *BaseModule[T]) IsRunning() bool {
	return b.State() == StateRunning
}

//--------------------------------
//		状态变化
//--------------------------------

// DoStart start时指定
func (b *BaseModule[T]) DoStart() {
	b.state.Store(uint32(StateRunning))
}

// DoClose close时执行
func (b *BaseModule[T]) DoClose() {
	b.state.Store(uint32(StateClosed))
}

// DoInit close时执行,允许重启
func (b *BaseModule[T]) DoInit() {
	b.state.Store(uint32(StateInit))
}

// CheckStart 是否可以start
func (b *BaseModule[T]) CheckStart() error {
	switch b.State() {
	case StateRunning:
		return ErrAlreadyStarted
	case StateClosed:
		return ErrAlreadyClosed
	default:
		return nil
	}
}

// CheckClose 是否可以close
func (b *BaseModule[T]) CheckClose() error {
	if b.State() == StateClosed {
		return ErrAlreadyClosed
	}
	return nil
}
