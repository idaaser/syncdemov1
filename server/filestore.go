package server

import (
	"encoding/json"
	"os"
	"sync"
)

func newJSONFileStore[T any](f string) *jsonFS[T] {
	return &jsonFS[T]{
		file: f,
		once: &sync.Once{},
	}
}

type jsonFS[T any] struct {
	file string

	once *sync.Once
	data []T
}

func (s *jsonFS[T]) load() []T {
	s.once.Do(func() {
		if content, err := os.ReadFile(s.file); err == nil {
			data := []T{}
			if err := json.Unmarshal(content, &data); err == nil {
				s.data = data
			}
		}
	})

	return s.data
}

func (s *jsonFS[T]) sublist(start, size int) ([]T, int) {
	return sublist(s.load(), start, size)
}

func sublist[T any](s []T, start, size int) ([]T, int) {
	l := len(s)
	if start >= l {
		return nil, -1
	}

	end := min(start+size, l)
	data := s[start:end]
	if end < l {
		return data, end
	}

	return data, -1
}

func min(i, j int) int {
	if i > j {
		return j
	}
	return i
}
