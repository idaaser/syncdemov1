package server

import (
	"encoding/json"
	"os"
	"sync"
)

func newJsonFileStore[T any](f string) *jsonFileStore[T] {
	return &jsonFileStore[T]{
		file: f,
		once: &sync.Once{},
	}
}

type jsonFileStore[T any] struct {
	file string

	once *sync.Once
	data []T
}

func (s *jsonFileStore[T]) load() []T {
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
