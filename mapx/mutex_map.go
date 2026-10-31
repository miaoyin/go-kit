package mapx

import "sync"

// RWMap 线程安全的map
//  1. sync.map缺点是类型断言, any类型判断
//  2. RWMap无需类型断言
type RWMap[K comparable, V any] struct {
	mu      sync.RWMutex
	dataMap map[K]V
}

func NewRWMap[K comparable, V any]() *RWMap[K, V] {
	return &RWMap[K, V]{
		dataMap: make(map[K]V),
	}
}

func (m *RWMap[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.dataMap[k]
	return v, ok
}

func (m *RWMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dataMap[k] = v
}

func (m *RWMap[K, V]) Delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.dataMap, k)
}

func (m *RWMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.dataMap)
}

func (m *RWMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]K, 0, len(m.dataMap))
	for k := range m.dataMap {
		keys = append(keys, k)
	}
	return keys
}

func (m *RWMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vals := make([]V, 0, len(m.dataMap))
	for _, v := range m.dataMap {
		vals = append(vals, v)
	}
	return vals
}

func (m *RWMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dataMap = make(map[K]V)
}

func (m *RWMap[K, V]) Range(fn func(k K, v V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.dataMap {
		if !fn(k, v) {
			break
		}
	}
}

func (m *RWMap[K, V]) RangeAndUpdate(fn func(k K, v V) (newV V, del bool)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.dataMap {
		newVal, del := fn(k, v)
		if del {
			delete(m.dataMap, k)
		} else {
			m.dataMap[k] = newVal
		}
	}
}

// GetOrSet 获取, 不存在则设置默认值
func (m *RWMap[K, V]) GetOrSet(k K, defaultVal V) V {
	m.mu.RLock()
	v, ok := m.dataMap[k]
	m.mu.RUnlock()
	if ok {
		return v
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.dataMap[k]; ok {
		return v
	}
	m.dataMap[k] = defaultVal
	return defaultVal
}

// GetAndDelete 读取后删除
func (m *RWMap[K, V]) GetAndDelete(k K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.dataMap[k]
	if ok {
		delete(m.dataMap, k)
	}
	return v, ok
}

func (m *RWMap[K, V]) SetIfNotExist(k K, v V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.dataMap[k]; ok {
		return false
	}
	m.dataMap[k] = v
	return true
}

func (m *RWMap[K, V]) Has(k K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.dataMap[k]
	return ok
}

func (m *RWMap[K, V]) SetMany(items map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range items {
		m.dataMap[k] = v
	}
}

func (m *RWMap[K, V]) DeleteMany(keys []K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		delete(m.dataMap, k)
	}
}

// Compute 根据k, 处理旧值
// 可减少大量零散操作
func (m *RWMap[K, V]) Compute(k K, fn func(old V, exist bool) (newV V, del bool)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	old, exist := m.dataMap[k]
	newV, del := fn(old, exist)
	if del {
		delete(m.dataMap, k)
	} else {
		m.dataMap[k] = newV
	}
}

func (m *RWMap[K, V]) Clone() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[K]V, len(m.dataMap))
	for k, v := range m.dataMap {
		out[k] = v
	}
	return out
}
