package confutil

import (
    "context"
    "os"
    "path/filepath"

    "github.com/fsnotify/fsnotify"
)

type ResourceEvent struct {
    SourceLayer string // base / remote / runtime
    ResourceID  string // 文件路径 / etcd key / badger key
    Data        []byte
    Err         error
}

//----------------------------------------------
//  文件
//----------------------------------------------

type FileWatcher struct{
    fsw *fsnotify.Watcher
}

func NewFileWatcher(fsw *fsnotify.Watcher) *FileWatcher {
    return &FileWatcher{fsw: fsw}
}

func (w *FileWatcher) WatchOne(ctx context.Context, filePath string) <-chan ResourceEvent {
    ch := make(chan ResourceEvent, 1)
    go func() {
        defer close(ch)
        abs, err := filepath.Abs(filePath)
        if err != nil {
            ch <- ResourceEvent{SourceLayer: "base", ResourceID: filePath, Err: err}
            return
        }

        dir := filepath.Dir(abs)
        name := filepath.Base(abs)
        //监控目录
        if err = w.fsw.Add(dir); err != nil {
            ch <- ResourceEvent{SourceLayer: "base", ResourceID: filePath, Err: err}
            return
        }

        for{
            select {
            case <-ctx.Done():
                return
            case ev, ok := <-w.fsw.Events:
                if !ok {
                    return
                }
                if filepath.Base(ev.Name) != name {
                    continue
                }
                //写入,新建,改名
                if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
                    continue
                }
                data, readErr := os.ReadFile(abs)
                ch <- ResourceEvent{
                    SourceLayer: "base",
                    ResourceID:  filePath,
                    Data:        data,
                    Err:         readErr,
                }
            case fswErr, ok := <-w.fsw.Errors:
                if !ok {
                    return
                }
                ch <- ResourceEvent{SourceLayer: "base", ResourceID: filePath, Err: fswErr}
                return
            }
        }
    }()
    return ch
}

