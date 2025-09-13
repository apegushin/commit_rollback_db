package db

import (
	"iter"
	"maps"
)

// TODO: reimplement Set using generics.
// Also look up if an interface exists for a Set data structure
// that is more widely used and implement to the interface

type Set struct {
	items map[int]struct{}
}

func NewSet() *Set {
	return &Set{
		items: make(map[int]struct{}),
	}
}

func (s *Set) Add(items ...int) {
	for _, item := range items {
		if !s.Contains(item) {
			s.items[item] = struct{}{}
		}
	}
}

func (s *Set) Contains(item int) bool {
	_, ok := s.items[item]
	return ok
}

func (s *Set) Items() iter.Seq[int] {
	return maps.Keys(s.items)
}

func (s *Set) Len() int {
	return len(s.items)
}

func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set) Remove(item int) {
	delete(s.items, item)
}

func (s *Set) Clear() {
	clear(s.items)
}
