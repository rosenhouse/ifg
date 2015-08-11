package application

import "fmt"

type MemoryDataStore struct {
	hash map[string][]byte
}

func NewMemoryDataStore() *MemoryDataStore { return &MemoryDataStore{make(map[string][]byte)} }

func (s *MemoryDataStore) Get(key string) ([]byte, error) {
	val, ok := s.hash[key]
	if !ok {
		return nil, fmt.Errorf("no value for key '%s'", key)
	}
	return val, nil
}

func (s *MemoryDataStore) Set(key string, val []byte) error {
	s.hash[key] = val
	return nil
}
