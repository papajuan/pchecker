package pchecker

/**
 * @author  papajuan
 * @date    10/10/2024
 **/

type SafeSet[K comparable] struct {
	SafeMap[K, bool]
}

func NewSafeSet[K comparable](m map[K]bool) *SafeSet[K] {
	return &SafeSet[K]{*NewSafeMap[K, bool](m)}
}

func (s *SafeSet[K]) Contains(key K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.m[key]
}

func (s *SafeSet[K]) Range(f func(K)) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.m {
		f(k)
	}
}
