package fileutil

import (
	"encoding/json"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

/*
	viper: support map, not support list
	encoding: supportmap, list
*/

func FileUnmarshal(path string, v any) error{
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		return yaml.Unmarshal(dataBytes, v)
	case ".json":
		return json.Unmarshal(dataBytes, v)
	case ".toml":
		return toml.Unmarshal(dataBytes, v)
	default:
		return json.Unmarshal(dataBytes, v)
	}
}

func FileMarshal(path string, v any) (err error){
	var dataBytes []byte
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		dataBytes, err = yaml.Marshal(v)
		if err != nil {
			return err
		}
	case ".json":
		dataBytes, err = json.Marshal(v)
		if err != nil {
			return err
		}
	case ".toml":
		dataBytes, err = toml.Marshal(v)
		if err != nil {
			return err
		}
	default:
		dataBytes, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(path, dataBytes, os.ModePerm)
}

func FileMarshalIndent(path string, v any) (err error){
	var dataBytes []byte
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		dataBytes, err = yaml.Marshal(v)
		if err != nil {
			return err
		}
	case ".json":
		dataBytes, err = json.MarshalIndent(v, "", "    ")
		if err != nil {
			return err
		}
	case ".toml":
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

