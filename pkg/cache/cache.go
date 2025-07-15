package cache

import (
	"iter"
	"maps"
	"strings"
	"sync"
)

type Storage struct {
	data map[string]any
	mu   sync.RWMutex
}

func New() *Storage {
	return &Storage{
		data: make(map[string]any),
	}
}

func (s *Storage) Set(k string, v any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k] = v
}

func (s *Storage) Get(k string) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.data[k]; ok {
		return v
	}
	return nil
}

func (s *Storage) Del(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, k)
}

func (s *Storage) DeleteByPrefix(pre string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k := range s.data {
		if strings.HasPrefix(k, pre) {
			delete(s.data, k)
		}
	}
}

func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]any)
}

func (s *Storage) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Storage) Keys() iter.Seq[string] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return maps.Keys(s.data)
}
