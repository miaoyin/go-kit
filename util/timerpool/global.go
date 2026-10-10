package timerpool

import (
	"sync"
)

var (
	//Default 避免到处传实例
	Default *TimerPool
	once    sync.Once
)

func InitDefault() {
	once.Do(func() {
		Default = NewTimerPool()
	})
}
