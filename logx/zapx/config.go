package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogSinkConfig struct {
	//文件选项
	FileOptions
	//zap选项
	ZapOptions
	//标识：main / api / error / access
	Name       string
	Console    bool
	// json / console
	Encoder string
	// 是否抄送error_sink
	CopyToErr bool
}

//FileOptions 文件选项
type FileOptions struct {
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

//ZapOptions 选项
type ZapOptions struct {
	//日志级别
	Level         string
	AddCaller     bool
	CallerSkip    int
	//AddStacktrace 打印栈信息
	AddStacktrace zapcore.Level
	EnableHooks []string
}

// MakeZapOptions 构造选项
func MakeZapOptions(cfg ZapOptions) []zap.Option {
	var opts []zap.Option

	if cfg.AddCaller {
		opts = append(opts, zap.AddCaller())
	}

	if cfg.AddStacktrace != zapcore.Level(-1) {
		opts = append(opts, zap.AddStacktrace(cfg.AddStacktrace))
	}

	if cfg.CallerSkip > 0 {
		opts = append(opts, zap.AddCallerSkip(cfg.CallerSkip))
	}

	var hooks []func(zapcore.Entry) error
	for _, name := range cfg.EnableHooks {
		if h, ok := hookRegistry[name]; ok {
			hooks = append(hooks, h)
		}
	}
	if len(hooks) > 0 {
		opts = append(opts, zap.Hooks(hooks...))
	}
	return opts
}
