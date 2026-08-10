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

//nopLogger 空实现
type nopLogger struct{}

func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}


//----------------------------------
//  全局日志(会产生隐式依赖)
// 1. 简单场景使用
// 2. 默认logger使用
//----------------------------------

var (
	globalLogger Logger = nopLogger{}
	mu     sync.RWMutex
)

func SetLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	if l == nil {
		//关闭日志
		globalLogger = nopLogger{}
		return
	}
	globalLogger = l
}

func GetLogger() Logger{
	mu.RLock()
	defer mu.RUnlock()
	return globalLogger
}
