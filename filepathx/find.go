package filepathx

import (
	"errors"
	"os"
)

// FindFirstFile 找文件, 找到就返回
func FindFirstFile(paths ...string) (string, error) {
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		if !info.IsDir() {
			return p, nil
		}
	}
	return "", errors.New("no valid file found")
}
