package logger

import "go.uber.org/zap"

type ZapAdapter struct {
	Logger *zap.SugaredLogger
}

func (z *ZapAdapter) Debug(msg string, args ...any) { z.Logger.Debugw(msg, args...) }
func (z *ZapAdapter) Info(msg string, args ...any)  { z.Logger.Infow(msg, args...) }
func (z *ZapAdapter) Warn(msg string, args ...any)  { z.Logger.Warnw(msg, args...) }
func (z *ZapAdapter) Error(msg string, args ...any) { z.Logger.Errorw(msg, args...) }