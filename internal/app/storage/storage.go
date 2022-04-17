package storage

import (
	"encoding/csv"
	"fmt"
	"os"
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

func ConstructStorage(fileStoragePath string) *Storage {
	s := &Storage{make(map[string]string), sync.RWMutex{}}

	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		fmt.Errorf("OpenFile error; %s", err)
		return s
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Errorf("ReadAll error; %s", err)
		return s
	}

	for _, rec := range records {
		s.URLsMap[rec[0]] = rec[1]
	}

	return s
}

func DestructStorage(fileStoragePath string, s *Storage) error {
	file, err := os.OpenFile(fileStoragePath, os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("OpenFile error; %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	var records [][]string
	for ID, URL := range s.URLsMap {
		records = append(records, []string{ID, URL})
	}

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("WriteAll error; %s", err)
	}

	writer.Flush()

	return nil
}
