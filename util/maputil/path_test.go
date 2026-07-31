package maputil

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
    m := map[string]any{
        "a": 1,
        "b": "2",
        "c": map[string]any{
            "c1": 3,
        },
    }
    require.Equal(t, GetByPathGeneric[int](m, "a"), 1)
    require.Equal(t, GetByPathGeneric[string](m, "b"), "2")
    require.Equal(t, GetByPathGeneric[map[string]any](m, "c"), map[string]any{"c1": 3})
}