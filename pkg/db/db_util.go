package db

import (
	"iter"
	"maps"
)

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		if !s.Contains(item) {
			s.items[item] = struct{}{}
		}
	}
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.items[item]
	return ok
}

func (s *Set[T]) Items() iter.Seq[T] {
	return maps.Keys(s.items)
}

func (s *Set[T]) Len() int {
	return len(s.items)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

func (s *Set[T]) Clear() {
	clear(s.items)
}
