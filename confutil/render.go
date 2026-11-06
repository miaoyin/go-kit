package confutil

import (
	"os"
	"regexp"
	"strings"
)

// varRegex 限定变量格式${xxx}
var varRegex = regexp.MustCompile(`\$\{([a-zA-Z0-9_\-]+)\}`)

type Lookup func(s string) string

// Render 变量扩展
// 1. 忽略大小写
// 2. 只支持${xxx}  ---os.Expand支持$x, 容易出问题
// 3. 未找到返回原始内容${xxx}  ---os.Expand会返回空
func Render(s string, varReader StringReader) string {

	return varRegex.ReplaceAllStringFunc(s, func(m string) string {
		sub := varRegex.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		rawKey := sub[1]
		normKey := strings.ToLower(rawKey)
		if v := varReader.GetString(normKey); v != "" {
			return v
		}

		//原始key
		if v := os.Getenv(rawKey); v != "" {
			return v
		}
		//大写
		if v := os.Getenv(strings.ToUpper(normKey)); v != "" {
			return v
		}
		//原样返回
		return m
	})
}

func RenderSlice(in []string, varReader StringReader) []string {
	out := make([]string, 0, len(in))
	for _, it := range in {
		out = append(out, Render(it, varReader))
	}
	return out
}
