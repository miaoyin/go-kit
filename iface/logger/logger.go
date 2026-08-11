package logger

import "sync"

//Logger 日志接口
// 支持zap
// go1.21+用log/slog
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

//NopLogger 空实现
type NopLogger struct{}

func (NopLogger) Debug(string, ...any) {}
func (NopLogger) Info(string, ...any)  {}
func (NopLogger) Warn(string, ...any)  {}
func (NopLogger) Error(string, ...any) {}


//----------------------------------
//  全局日志(会产生隐式依赖)
// 1. 简单场景使用
// 2. 默认logger使用
//----------------------------------

var (
	globalLogger Logger = NopLogger{}
	mu     sync.RWMutex
)

func SetLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	if l == nil {
		//关闭日志
		globalLogger = NopLogger{}
		return
	}
	globalLogger = l
}

func GetLogger() Logger{
	mu.RLock()
	defer mu.RUnlock()
	return globalLogger
}
