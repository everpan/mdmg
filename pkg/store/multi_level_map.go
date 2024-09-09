package store

import (
	"sync"
)

type OneLevelMap[K comparable, V any] struct {
	m     map[K]V
	rwMux sync.RWMutex
}
type TwoLevelMap[K1 comparable, K2 comparable, V any] struct {
	m     map[K1]OneLevelMap[K2, V]
	rwMux sync.RWMutex
}

func (tlm TwoLevelMap[K1, K2, V]) Acquire(k1 K1) OneLevelMap[K2, V] {
	tlm.rwMux.Lock()
	defer tlm.rwMux.Unlock()
	if tlm.m == nil {
		tlm.m = make(map[K1]OneLevelMap[K2, V])
	}
	m, ok := tlm.m[k1]
	if !ok {
		m = OneLevelMap[K2, V]{}
		m.m = make(map[K2]V)
		tlm.m[k1] = m
	}
	return m
}

func (m OneLevelMap[K, V]) Set(k K, v V) {
	m.rwMux.Lock()
	defer m.rwMux.Unlock()
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[k] = v
}

func (m OneLevelMap[K, V]) Release(k K) {
	if m.m == nil {
		return
	}
	m.rwMux.Lock()
	defer m.rwMux.Unlock()
	delete(m.m, k)
}

func (m OneLevelMap[K, V]) Get(k K) (v V, ok bool) {
	m.rwMux.RLock()
	defer m.rwMux.RUnlock()
	if m.m == nil {
		return v, false
	}
	v, ok = m.m[k]
	return v, ok
}
