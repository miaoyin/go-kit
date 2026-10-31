package slicex

import (
	"fmt"
	"testing"
)

func TestDiffString(t *testing.T) {
	d := Difference[string]([]string{"a", "b"}, []string{"c", "d"})
	fmt.Println(d)
}
