package confutil

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestRender_base(t *testing.T) {
    dataItems := []struct{
        Path string
        VarStores map[string]string
        Expected string
        Desc    string
    }{
        {
            Path: "${home}/abc.toml",
            VarStores: map[string]string{
                "home": "/config",
            },
            Expected: "/config/abc.toml",
            Desc: "变量替换",
        },
        {
            Path: "${Home}/${name}/abc.toml",
            VarStores: map[string]string{
                "home": "/config",
                "NAME": "name",
            },
            Expected: "/config/name/abc.toml",
            Desc: "忽略大小写",
        },
        {
            Path: "${home}/abc.toml",
            VarStores: map[string]string{
            },
            Expected: "${home}/abc.toml",
            Desc: "未找到变量, 保持不变",
        },
        {
            Path: "$home/abc.toml",
            VarStores: map[string]string{
                "home": "/config",
            },
            Expected: "$home/abc.toml",
            Desc: "不支持$x",
        },
    }

    for _, dataItem := range dataItems {
        val := Render(dataItem.Path, dataItem.VarStores)
        t.Log(dataItem.Desc, val)
        assert.Equal(t, dataItem.Expected, val)
    }
}