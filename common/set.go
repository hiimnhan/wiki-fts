package common

import "sync"

type Set struct {
	items map[any]bool
	lock  sync.RWMutex
}

func (s *Set) Add(item any) *Set {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.items == nil {
		s.items = make(map[any]bool)
	}

	if _, ok := s.items[item]; !ok {
		s.items[item] = true
	}

	return s
}

func (s *Set) Delete(item any) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.items[item]
	if ok {
		delete(s.items, item)
	}

	return ok
}

func (s *Set) Has(item any) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.items[item]
	return ok
}

func (s *Set) Items() []any {
	s.lock.RLock()
	defer s.lock.RUnlock()
	items := []any{}
	for i := range s.items {
		items = append(items, i)
	}

	return items
}

func (s *Set) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
