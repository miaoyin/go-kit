package osutil

import (
    "os"
)


func Executable() string{
    exePath, _ := os.Executable()
    return exePath
}

func IsFile(path string) bool {
    f, e := os.Stat(path)
    if e != nil {
        return false
    }
    return !f.IsDir()
}

func IsDir(dir string) bool {
    f, e := os.Stat(dir)
    if e != nil {
        return false
    }
    return f.IsDir()
}

// IsExist 是否存在, 返回bool
func IsExist(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

// IsExistE 是否存在, 返回error
func IsExistE(path string) error {
    _, err := os.Stat(path)
    if err ==nil || os.IsExist(err){
        return nil
    }
    return err
}
