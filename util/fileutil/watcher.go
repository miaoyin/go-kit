package fileutil

import (
    "sync"

    "github.com/fsnotify/fsnotify"
)

type Subscriber struct {
    Event   chan fsnotify.Event
    Done chan struct{}
}

type DirectoryWatcher struct{
    *fsnotify.Watcher
    subs    map[*Subscriber]struct{}
    rmu sync.RWMutex
}

func NewDirectoryWatcher() *DirectoryWatcher {
    fsWatcher, _ := fsnotify.NewWatcher()
    b := &DirectoryWatcher{
        Watcher: fsWatcher,
        subs: map[*Subscriber]struct{}{},
    }
    go b.dispatchEvent()
    return b
}

func (b *DirectoryWatcher) dispatchEvent() {
    for{
        select{
        case evt, ok := <-b.Watcher.Events:
            if !ok {
                return
            }
            b.broadcast(evt)
        case _ = <-b.Watcher.Errors:
            return
        }
    }
}

func (b *DirectoryWatcher) broadcast(evt fsnotify.Event) {
    b.rmu.RLock()
    defer b.rmu.RUnlock()
    for sub, _ := range b.subs {
        select {
        case <-sub.Done:
        case sub.Event <- evt:
        default:
        }
    }
}

func (b *DirectoryWatcher) Subscribe() *Subscriber {
    sub := &Subscriber{
        Event:   make(chan fsnotify.Event),
        Done: make(chan struct{}),
    }
    b.rmu.Lock()
    defer b.rmu.Unlock()
    b.subs[sub] = struct{}{}
    return sub
}

func (b *DirectoryWatcher) Unsubscribe(sub *Subscriber) {
    b.rmu.Lock()
    defer b.rmu.Unlock()

    delete(b.subs, sub)
    close(sub.Done)
}
