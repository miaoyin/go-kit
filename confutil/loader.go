package confutil

import (
    "bytes"
    "context"
    "io"
    "io/fs"
    "os"
    "path"

    "github.com/miaoyin/go-kit/util/fileutil"
)

//FileLoader 加载文件
type FileLoader struct{
    fileSys fs.FS
    vars map[string]string
}

func NewFileLoader(fileSys fs.FS, vars map[string]string) *FileLoader {
    return &FileLoader{fileSys: fileSys, vars: vars}
}

func (l *FileLoader) LoadOne(_ context.Context, filePath string) (map[string]any, error) {
    // 1.打开文件
    filePath = Render(filePath, l.vars)
    f, err := l.fileSys.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    // 2.读入buf
    var buf bytes.Buffer
    _, err = io.Copy(&buf, f)
    if err != nil {
        return nil, err
    }
    raw := buf.Bytes()
    rawText := Render(string(raw), l.vars)

    // 3.转为map
    var rawMap map[string]any
    err = fileutil.Unmarshal(path.Ext(filePath), []byte(rawText), &rawMap)
    if err != nil {
        return nil, err
    }
    return rawMap, nil
}

//AbsFS 支持绝对路径
type AbsFS struct{}

func (a AbsFS) Open(name string) (fs.File, error) {
    return os.Open(name)
}
