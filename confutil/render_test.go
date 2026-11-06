package confutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender_base(t *testing.T) {
	dataItems := []struct {
		Path     string
		ReadFunc StringReaderFunc
		Expected string
		Desc     string
	}{
		{
			Path: "${home}/abc.toml",
			ReadFunc: func(key string) string {
				return ""
			},
			Expected: "/config/abc.toml",
			Desc:     "变量替换",
		},
		{
			Path: "${Home}/${name}/abc.toml",
			ReadFunc: func(key string) string {
				dataMap := map[string]string{
					"home": "/config",
					"NAME": "name",
				}
				return dataMap[key]
			},
			Expected: "/config/name/abc.toml",
			Desc:     "忽略大小写",
		},
		{
			Path: "${home}/abc.toml",
			ReadFunc: func(key string) string {
				return ""
			},
			Expected: "${home}/abc.toml",
			Desc:     "未找到变量, 保持不变",
		},
		{
			Path: "$home/abc.toml",
			ReadFunc: func(key string) string {
				dataMap := map[string]string{
					"home": "/config",
				}
				return dataMap[key]
			},
			Expected: "$home/abc.toml",
			Desc:     "不支持$x",
		},
	}

	for _, dataItem := range dataItems {
		val := Render(dataItem.Path, dataItem.ReadFunc)
		t.Log(dataItem.Desc, val)
		assert.Equal(t, dataItem.Expected, val)
	}
}
