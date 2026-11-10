package confutil

import (
	"github.com/miaoyin/go-kit/atomicx"
	"github.com/spf13/viper"
)

// ConfigReader 读接口, 使用者
type ConfigReader interface {
	GetString(key string, def ...string) string
	GetInt(key string, def ...int) int
	GetBool(key string, def ...bool) bool
	GetFloat64(key string, def ...float64) float64
	Sub(prefix string) ConfigReader
	IsSet(key string) bool
	AllSettings() map[string]any
	Unmarshal(raw any) error
}

// Reloader 加载, 支持热更新
type Reloader interface {
	Swap(newVP *viper.Viper) *viper.Viper
}

type configReader struct {
	raw *atomicx.PointerHolder[viper.Viper]
}

// NewConfigReader 返回两个接口
// 使用者不要看到Reloader接口
func NewConfigReader(raw *atomicx.PointerHolder[viper.Viper]) ConfigReader {
	return &configReader{raw: raw}
}

func (d *configReader) snap() *viper.Viper {
	v, err := d.raw.GetSnapshot()
	if err != nil {
		return viper.New()
	}
	return v
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
	return NewConfigReader(atomicx.NewPointerHolder(subVP))
}

func (d *configReader) IsSet(key string) bool {
	return d.snap().IsSet(key)
}

func (d *configReader) AllSettings() map[string]any {
	return d.snap().AllSettings()
}

func (d *configReader) Unmarshal(rawVal any) error {
	return d.snap().Unmarshal(rawVal)
}
