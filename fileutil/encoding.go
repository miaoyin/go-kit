package fileutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// viper: support map, not support list
// encoding: support map, list


//Unmarshal 解析
// ext后缀
func Unmarshal(ext string, src []byte, v any) error {
	ext = strings.TrimLeft(ext, ".")
	switch ext {
	case "yaml", "yml":
		return yaml.Unmarshal(src, v)
	case "json":
		return json.Unmarshal(src, v)
	case "toml":
		return toml.Unmarshal(src, v)
	default:
		return fmt.Errorf("unknown extension: %s", ext)
	}
}

func Marshal(ext string, v any) ([]byte, error) {
	ext = strings.TrimLeft(ext, ".")
	switch ext {
	case "yaml", "yml":
		return yaml.Marshal(v)
	case "json":
		return json.Marshal(v)
	case "toml":
		return toml.Marshal(v)
	default:
		return nil, fmt.Errorf("unknown extension: %s", ext)
	}
}

// FileUnmarshal 读文件并解码
func FileUnmarshal(path string, v any) error {
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return Unmarshal(filepath.Ext(path), dataBytes, v)
}

// MarshalByPath 编码
//    path 文件路径或者后缀
//    v 数据对象
func MarshalByPath(path string, v any) ([]byte, error) {
	return Marshal(filepath.Ext(path), v)
}

// FileMarshal 编码并保存文件
func FileMarshal(path string, v any) (err error) {
    dataBytes, err := MarshalByPath(path, v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, dataBytes, os.ModePerm)
}

func FileMarshalIndent(path string, v any) (err error) {
	var dataBytes []byte
    ext := strings.TrimLeft(filepath.Ext(path), ".")
	switch ext {
	case "yaml", "yml":
		dataBytes, err = yaml.Marshal(v)
		if err != nil {
			return err
		}
	case "json":
		dataBytes, err = json.MarshalIndent(v, "", "    ")
		if err != nil {
			return err
		}
	case "toml":
		dataBytes, err = toml.Marshal(v)
		if err != nil {
			return err
		}
	default:
		dataBytes, err = json.MarshalIndent(v, "", "    ")
		if err != nil {
			return err
		}
	}
	return os.WriteFile(path, dataBytes, os.ModePerm)
}
