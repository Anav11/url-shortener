package storage

import (
	"fmt"
	"sync"
)

type Repository interface {
	Add(ID string, URL string) error
	Get(ID string) (string, error)
}

type Storage struct {
	URLsMap map[string]string
	mutex   sync.RWMutex
}

func (s *Storage) Add(ID string, URL string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.URLsMap[ID] != "" {
		return fmt.Errorf(`ID=%s; URL already exists`, ID)
	}

	s.URLsMap[ID] = URL

	return nil
}

func (s *Storage) Get(ID string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	URL := s.URLsMap[ID]
	if URL == "" {
		return "", fmt.Errorf("URL not found")
	}

	return URL, nil
}

func ConstructStorage() *Storage {
	return &Storage{make(map[string]string), sync.RWMutex{}}
}
