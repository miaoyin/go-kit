package zapx

import "go.uber.org/zap/zapcore"

//------------------------------------------
//   hook钩子
//------------------------------------------

var hookRegistry = make(map[string]func(zapcore.Entry) error)

// RegisterHook 注册钩子
func RegisterHook(name string, h func(zapcore.Entry) error) {
    hookRegistry[name] = h
}

func GetHook(name string) func(zapcore.Entry) error {
    v, ok := hookRegistry[name]
    if !ok {
        return nil
    }
    return v
}


//------------------------------------------
//   encoder编码
//------------------------------------------

var encoderRegistry = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
    "json":    zapcore.NewJSONEncoder,
    "console": zapcore.NewConsoleEncoder,
}

// RegisterEncoder 注册钩子
func RegisterEncoder(name string, e func(zapcore.EncoderConfig) zapcore.Encoder) {
    encoderRegistry[name] = e
}

func GetEncoder(name string)func(zapcore.EncoderConfig) zapcore.Encoder{
    v, ok := encoderRegistry[name]
    if !ok {
        return nil
    }
    return v
}
