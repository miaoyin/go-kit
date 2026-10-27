package confutil

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestMerge_base(t *testing.T) {
    dataItems := []struct{
        Dst  map[string]interface{}
        Src  map[string]interface{}
        ArrayMode MergeMode
        Expected map[string]interface{}
        Desc string
    }{
        {
            Dst:  map[string]interface{}{
                "files": []string{"1.toml", "2.toml"},
            },
            Src:  map[string]interface{}{
                "files": []string{"3.toml", "4.toml"},
            },
            ArrayMode: ArrayAppend,
            Expected: map[string]interface{}{
                "files": []string{"1.toml", "2.toml", "3.toml", "4.toml"},
            },
            Desc: "测试ArrayAppend",
        },

        {
            Dst:  map[string]interface{}{
                "files": []string{"1.toml", "2.toml"},
            },
            Src:  map[string]interface{}{
                "files": []string{"3.toml", "4.toml"},
            },
            ArrayMode: ArrayReplace,
            Expected: map[string]interface{}{
                "files": []string{"3.toml", "4.toml"},
            },
            Desc: "测试ArrayReplace",
        },
    }
    for _, dataItem := range dataItems {
        err := Merge(dataItem.Dst, dataItem.Src, dataItem.ArrayMode)
        assert.NoError(t, err)
        t.Log(dataItem.Desc, dataItem.Dst, dataItem.ArrayMode)
        assert.Equal(t, dataItem.Expected, dataItem.Dst)
    }
}