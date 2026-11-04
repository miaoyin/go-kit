package zapx

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type SinkWriter struct {
	Name    string
	WS      zapcore.WriteSyncer
	Closer  *SafeLumberjack
	Encoder zapcore.Encoder
}

// NewSinkWriter 根据sink配置构建writer+encoder，启动阶段调用一次
func NewSinkWriter(cfg LogSinkConfig) (*SinkWriter, error) {
	// encoder工厂，可注册扩展
	var enc zapcore.Encoder
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	if fn, ok := encoderRegistry[cfg.Encoder]; !ok {
		return nil, fmt.Errorf("unknown encoder %s", cfg.Encoder)
	} else {
		enc = fn(encCfg)
	}

	var ws zapcore.WriteSyncer
	var closer *SafeLumberjack
	if cfg.FilePath != "" {
		lj := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		safe := NewSafeLumberjack(lj)
		ws = zapcore.AddSync(safe)
		closer = safe
	}
	//是否写控制台
	if cfg.Console {
		if ws == nil {
			ws = zapcore.AddSync(os.Stdout)
		} else {
			ws = zapcore.NewMultiWriteSyncer(ws, zapcore.AddSync(os.Stdout))
		}
	}
	//默认
	if ws == nil {
		ws = zapcore.AddSync(os.Stdout)
	}

	return &SinkWriter{
		Name:    cfg.Name,
		WS:      ws,
		Closer:  closer,
		Encoder: enc,
	}, nil
}

// BuildCore 创建core
// 如果要level分流,可以自定义TeeErrorCore
func BuildCore(lvl zapcore.Level, sink *SinkWriter) zapcore.Core {
	return zapcore.NewCore(sink.Encoder, sink.WS, lvl)
}
