package osutil

import (
    "os"
)

/*---------------------------------

osutil

---------------------------------*/

func CurrentExePath() string{
    //name = filepath.Base(xx)
    //dir, name = filepath.Split(xx)
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

func IsExist(path string) bool {
    // 是否存在, 返回bool
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

func IsExistE(path string) error {
    //是否存在, 返回error
    _, err := os.Stat(path)
    if err ==nil || os.IsExist(err){
        return nil
    }
    return err
}
