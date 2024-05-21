package common

import "sync"

type Set struct {
	items map[int]bool
	lock  sync.RWMutex
}

func (s *Set) Add(item int) *Set {
	s.lock.Lock()

	if s.items == nil {
		s.items = make(map[int]bool)
	}
	s.lock.Unlock()

	s.lock.RLock()
	if _, ok := s.items[item]; !ok {
		s.items[item] = true
	}
	s.lock.RUnlock()

	return s
}

func (s *Set) Delete(item int) bool {
	s.lock.Lock()
	_, ok := s.items[item]
	if ok {
		delete(s.items, item)
	}
	s.lock.Unlock()

	return ok
}

func (s *Set) Merge(another *Set) *Set {
	s.lock.Lock()
	for item := range another.Items() {
		s.items[item] = true
	}
	s.lock.Unlock()

	return s
}

func (s *Set) Has(item int) bool {
	s.lock.RLock()
	_, ok := s.items[item]
	s.lock.RUnlock()
	return ok
}

func (s *Set) Items() []int {
	s.lock.RLock()
	items := []int{}
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
