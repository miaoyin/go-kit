package sliceutil

import (
    "fmt"
    "testing"
)

func TestDiffString(t *testing.T) {
    d := DiffSlice[string]([]string{"a", "b"}, []string{"c", "d"})
    fmt.Println(d)
}