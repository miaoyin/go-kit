package module

import (
	"context"
	"errors"
)


//State 模块状态
type State uint8

const (
	//StateInit 已创建，未启动
	StateInit State = iota
	//StateRunning 正常运行, running时持续状态, Started瞬态不要用started
	StateRunning
	//StateClosed 已关闭（终态，无法重启）
	StateClosed
)

func (s State) String() string {
	switch s {
	case StateInit:
		return "init"
	case StateRunning:
		return "running"
	case StateClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// 标准错误
var (
	ErrAlreadyStarted = errors.New("module already started")
	ErrNotRunning     = errors.New("module not running")
	ErrAlreadyClosed  = errors.New("module already closed")
	//ErrNeedRestart 字段不支持热更新时返回
	ErrNeedRestart    = errors.New("config changed non-hot-reload field, need restart module")
	ErrModuleDisabled = errors.New("module disabled by config")
)

//RawModule 通用模块接口
// 用于统一管理
type RawModule interface {
	//Name 全局唯一标识
	Name() string
	//Start 启动
	// 1. ctx用于启动过程中的阻塞操作, 一次性ctx, 不要存在struct中
	// 2. NewModule时已经初始过配置
	Start(ctx context.Context) error

	//Close 关闭
	Close(ctx context.Context) error

	//State 模块状态
	State() State

	//IsRunning 快捷接口
	IsRunning() bool
}

//Module 专用模块管理
// T是配置参数类型
// T推荐用struct类型,避免用指针. 配置参数热更新,使用副本更安全
type Module[T any] interface{
	//RawModule 通用接口
	RawModule

	//SetConfig 设置配置,运行时修改. 初始配置在New时传入
	// 1. Init状态：更新待生效配置
	// 2. Running状态：尝试热更新，存在不可热更字段返回 ErrNeedRestart
	// 3. Closed状态：直接返回错误
	SetConfig(cfg T) error

	//GetConfig 读取配置
	// 1.参数*T不安全, 使用时先获取副本
	// 2.int64是版本号, 区分参数是否修改
	GetConfig() (*T, int64)
}

