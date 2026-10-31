package confutil

import (
    "sync/atomic"
    "time"

    "github.com/spf13/viper"
)

//ConfigReader 读接口, 使用者
type ConfigReader interface{
    GetString(key string, def ...string) string
    GetInt(key string, def ...int) int
    GetBool(key string, def ...bool) bool
    GetFloat64(key string, def ...float64) float64
    Sub(prefix string) ConfigReader
    IsSet(key string) bool
    AllSettings() map[string]any
    Unmarshal(raw any) error
}

//Reloader 加载, 支持热更新
type Reloader interface {
    Reload(newVP *viper.Viper)
}

//viperSnapshot 快照信息
// 1. 预留, 可以增加附加信息
// 2. 减小后续改动
type viperSnapshot struct {
    vp *viper.Viper
    version   uint64
    changedAt time.Time
}

type configReader struct {
    raw atomic.Value
}

func NewConfigReader(initVP *viper.Viper) ConfigReader {
    d := &configReader{}
    d.raw.Store(&viperSnapshot{vp: initVP})
    return d
}

//NewConfigReaderWithReloader 返回两个接口
// 使用者不要看到Reloader接口
func NewConfigReaderWithReloader(initVP *viper.Viper) (ConfigReader, Reloader) {
    d := &configReader{}
    d.raw.Store(initVP)
    return d, d
}

//Reload 支持热更新
func (d *configReader) Reload(newVP *viper.Viper) {
    d.raw.Store(&viperSnapshot{vp: newVP})
}

func (d *configReader) snap() *viper.Viper {
    return d.raw.Load().(*viperSnapshot).vp
}

func (d *configReader) GetString(key string, def ...string) string {
    v := d.snap()
    if len(def) > 0 && !v.IsSet(key) {
        return def[0]
    }
    return v.GetString(key)
}

func (d *configReader) GetInt(key string, def ...int) int {
    v := d.snap()
    if len(def) > 0 && !v.IsSet(key) {
        return def[0]
    }
    return v.GetInt(key)
}

func (d *configReader) GetBool(key string, def ...bool) bool {
    v := d.snap()
    if len(def) > 0 && !v.IsSet(key) {
        return def[0]
    }
    return v.GetBool(key)
}

func (d *configReader) GetFloat64(key string, def ...float64) float64 {
    v := d.snap()
    if len(def) > 0 && !v.IsSet(key) {
        return def[0]
    }
    return v.GetFloat64(key)
}

func (d *configReader) Sub(prefix string) ConfigReader {
    subVP := d.snap().Sub(prefix)
    if subVP == nil {
        subVP = viper.New()
    }
    return NewConfigReader(subVP)
}

func (d *configReader) IsSet(key string) bool {
    return d.snap().IsSet(key)
}

func (d *configReader) AllSettings() map[string]any {
    return d.snap().AllSettings()
}

func (d *configReader) Unmarshal(rawVal any) error{
    return d.snap().Unmarshal(rawVal)
}
