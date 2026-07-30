package module

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//Manager 模块管理
// 适用:
//   1.批量Start, 批量Close, 状态查询
// 不适合:
//   1.模块之间通过Manager查找调用
type Manager struct{
	mu      sync.RWMutex
	//ordered 保留注册顺序
	ordered []RawModule
	//查找
	index   map[string]RawModule
}

func NewManager() *Manager {
	return &Manager{
		index: make(map[string]RawModule),
	}
}

func (m *Manager) Register(mod RawModule) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := mod.Name()
	if _, exists := m.index[name]; exists {
		return fmt.Errorf("module %s already registered", name)
	}
	m.index[name] = mod
	m.ordered = append(m.ordered, mod)
	return nil
}

// Get 根据名称获取模块
func (m *Manager) Get(name string) RawModule {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.index[name]
}

// List 返回所有模块
func (m *Manager) List() []RawModule {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]RawModule, len(m.ordered))
	copy(out, m.ordered)
	return out
}

//StartAll 批量启动, 支持回滚
func (m *Manager) StartAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	started := make([]RawModule, 0, len(m.ordered))
	for _, mod := range m.ordered {
		if err := mod.Start(ctx); err != nil {
			// 回滚
			closeCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			for i := len(started) - 1; i >= 0; i-- {
				_ = started[i].Close(closeCtx)
			}
			return fmt.Errorf("start module[%s] failed: %w", mod.Name(), err)
		}
		started = append(started, mod)
	}
	return nil
}

// CloseAll 反序关闭
func (m *Manager) CloseAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 逆序遍历
	for i := len(m.ordered) - 1; i >= 0; i-- {
		mod := m.ordered[i]
		if err := mod.Close(ctx); err != nil {
			return fmt.Errorf("close module[%s] failed: %w", mod.Name(), err)
		}
	}
	return nil
}