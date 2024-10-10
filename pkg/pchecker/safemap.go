package pchecker

import "sync"

/**
 * @author  papajuan
 * @date    10/10/2024
 **/

type SafeMap[K comparable, V any] struct {
	m    map[K]V
	lock sync.RWMutex
}

func NewSafeMap[K comparable, V any](m map[K]V) *SafeMap[K, V] {
	result := SafeMap[K, V]{}
	if m != nil && len(m) > 0 {
		result.m = m
	} else {
		result.m = make(map[K]V)
	}
	return &result
}

func (s *SafeMap[K, V]) Contains(key K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.m[key]
	return ok
}

func (s *SafeMap[K, V]) PutAll(mAdd map[K]V) {
	if mAdd != nil && len(mAdd) > 0 {
		s.lock.Lock()
		defer s.lock.Unlock()
		for k, v := range mAdd {
			s.m[k] = v
		}
	}
}

func (s *SafeMap[K, V]) Range(f func(K, V)) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k, v := range s.m {
		f(k, v)
	}
}
