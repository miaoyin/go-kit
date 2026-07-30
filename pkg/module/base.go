package module

import (
	"sync"
	"sync/atomic"
)

//EmptyConfig 无参数配置
// 例如 NewBaseModule[EmptyConfig]("name", EmptyConfig{})
type EmptyConfig struct{}

//BaseModule 基础模块
// 避免重复接口
type BaseModule[T any] struct{
	name    string
	//muState 状态使用锁, 保证业务一致
	muState sync.RWMutex
	state   State
	//confVal 配置, 支持热更新
	confVal atomic.Value
	//confVer 配置版本, 默认0
	confVer atomic.Int64
}

func NewBaseModule[T any](name string, initCfg T) *BaseModule[T] {
	b := &BaseModule[T]{
		name:  name,
		state: StateInit,
	}
	b.confVal.Store(&initCfg)
	return b
}

func (b *BaseModule[T]) Name() string {
	return b.name
}

//GetConfig 读取配置, 返回父辈
// 1.高频调用时, 使用缓存, 定时刷新减少调用频率
// 2.configVer版本号, 区分是否更新
func (b *BaseModule[T]) GetConfig() (*T, int64) {
	ptr := b.confVal.Load().(*T)
	//返回副本
	return ptr, b.confVer.Load()
}

//SetConfig 设置配置
// 根据需要覆盖
func (b *BaseModule[T]) SetConfig(cfg T) error{
	b.confVal.Store(&cfg)
	b.confVer.Add(1)
	return nil
}

func (b *BaseModule[T]) State() State {
	b.muState.RLock()
	defer b.muState.RUnlock()
	return b.state
}

func (b *BaseModule[T]) IsRunning() bool {
	return b.State() == StateRunning
}

//--------------------------------
//		状态变化
//--------------------------------

//DoStart start时指定
func (b *BaseModule[T]) DoStart() {
	b.muState.Lock()
	defer b.muState.Unlock()
	b.state = StateRunning
}

//DoClose close时执行
func (b *BaseModule[T]) DoClose() {
	b.muState.Lock()
	defer b.muState.Unlock()
	b.state = StateClosed
}

//CheckStart 是否可以start
func (b *BaseModule[T]) CheckStart() error {
	b.muState.RLock()
	defer b.muState.RUnlock()
	switch b.state {
	case StateRunning:
		return ErrAlreadyStarted
	case StateClosed:
		return ErrAlreadyClosed
	default:
		return nil
	}
}

//CheckClose 是否可以close
func (b *BaseModule[T]) CheckClose() error {
	b.muState.RLock()
	defer b.muState.RUnlock()
	if b.state == StateClosed {
		return ErrAlreadyClosed
	}
	return nil
}