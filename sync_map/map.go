package sync

import "sync"

type Map[K any, V any] struct {
	syncMap *sync.Map
}

func NewMap[K any, V any]() *Map[K, V] {
	syncMap := new(sync.Map)
	mpData := new(Map[K, V])
	mpData.syncMap = syncMap
	return mpData
}

// Set 设置数据
func (s *Map[K, V]) Set(key K, value V) {
	s.syncMap.Store(key, value)
	return
}

// Get 获取数据
func (s *Map[K, V]) Get(key K) (value V) {
	val, ok := s.syncMap.Load(key)
	if !ok {
		return
	}
	if val != nil {
		value, ok = val.(V)
		if !ok {
			return
		}
	}
	return
}

func (s *Map[K, V]) Delete(key K) {
	s.syncMap.Delete(key)
}
