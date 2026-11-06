package filepathx

import "path/filepath"

// ToLinuxPath 路径格式，否则变量替换会报错
// 格式 /name1/name2/name3
func ToLinuxPath(fp string) string {
	return filepath.ToSlash(filepath.Clean(fp))
}
