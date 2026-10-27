package confutil

import (
    "bytes"
    "io"
    "os"

    "github.com/mitchellh/mapstructure"
    "github.com/spf13/viper"
)

//ReadSettings 从r读取配置
func ReadSettings(r io.Reader, cfgType string, strict bool) (map[string]any, error) {
    v := viper.New()
    v.SetConfigType(cfgType)
    if err := v.ReadConfig(r); err != nil {
        return nil, err
    }

    if !strict {
        return v.AllSettings(), nil
    }

    // strict模式：用 UnmarshalExact 校验未知key
    var tmp map[string]any
    decOpt := viper.DecoderConfigOption(func(dc *mapstructure.DecoderConfig) {
        dc.ErrorUnused = true
    })
    if err := v.UnmarshalExact(&tmp, decOpt); err != nil {
        return nil, err
    }
    return tmp, nil
}

//ReadSettingsFile 读取配置文件
func ReadSettingsFile(filePath string, cfgType string, strict bool) (map[string]any, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    return ReadSettings(bytes.NewReader(data), cfgType, strict)
}
