package common

import "sync"

type Set struct {
	items map[any]bool
	lock  sync.RWMutex
}

func (s *Set) Add(item any) *Set {
	s.lock.Lock()

	if s.items == nil {
		s.items = make(map[any]bool)
	}
	s.lock.Unlock()

	s.lock.RLock()
	if _, ok := s.items[item]; !ok {
		s.items[item] = true
	}
	s.lock.RUnlock()

	return s
}

func (s *Set) Delete(item any) bool {
	s.lock.Lock()
	_, ok := s.items[item]
	if ok {
		delete(s.items, item)
	}
	s.lock.Unlock()

	return ok
}

func (s *Set) Has(item any) bool {
	s.lock.RLock()
	_, ok := s.items[item]
	s.lock.RUnlock()
	return ok
}

func (s *Set) Items() []any {
	s.lock.RLock()
	items := []any{}
	for i := range s.items {
		items = append(items, i)
	}
	s.lock.RUnlock()

	return items
}

func (s *Set) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
