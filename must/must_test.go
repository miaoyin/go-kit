package must

import (
	"fmt"
	"os"
	"testing"
)

func TestMust_get(t *testing.T) {
	bytes := Get(os.ReadFile("a.toml"))
	fmt.Println(bytes)
}

func TestMust_ok(t *testing.T) {
	dMap := map[string]interface{}{}
	v, ok := dMap["key"]
	vv := OK(v, ok)
	fmt.Println(vv)
}

func TestMust_err(t *testing.T) {
	NoErr(os.MkdirAll("./tmp", 0755))
}
