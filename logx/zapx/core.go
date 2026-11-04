package zapx

import (
    "go.uber.org/zap/zapcore"
)

//----------------------------------------------
//  日志name过滤
//----------------------------------------------

//FilterCore 过滤
type FilterCore struct {
    inner zapcore.Core
    wantName   string
}

func NewFilterCore(inner zapcore.Core, wantName string) zapcore.Core {
    return &FilterCore{
        inner: inner,
        wantName:   wantName,
    }
}

//Check 过滤
// apiLog = logger.Named("api")
func (c *FilterCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
    if entry.LoggerName != c.wantName {
        return nil
    }
    return c.inner.Check(entry, ce)
}

func (c *FilterCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
    return c.inner.Write(entry, fields)
}

func (c *FilterCore) Sync() error {
    return c.inner.Sync()
}

func (c *FilterCore) Enabled(lvl zapcore.Level) bool {
    return c.inner.Enabled(lvl)
}

func (c *FilterCore) With(fds []zapcore.Field) zapcore.Core {
    return c.inner.With(fds)
}
