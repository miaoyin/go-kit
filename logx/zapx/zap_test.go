package zapx

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func TestZap_base(t *testing.T) {
    apiSink := LogSinkConfig{
        Name: "api",
        FileOptions: FileOptions{
            FilePath: "d:/a.log",
        },
        Encoder: "console",
    }
    w, err  := NewSinkWriter(apiSink)
    assert.NoError(t, err)
    core := BuildCore(zapcore.InfoLevel, w)
    opts := MakeZapOptions(apiSink.ZapOptions)
    logger := zap.New(core, opts...)
    logger.Info("hello world")
}