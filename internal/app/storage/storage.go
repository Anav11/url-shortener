package storage

import "fmt"

type Repository interface {
	Add(ID string, URL string) error
	Get(ID string) (string, error)
}

type Storage struct {
	URLsMap map[string]string
}

func (s Storage) Add(ID string, URL string) error {
	if existingURL, _ := s.Get(ID); existingURL != "" {
		return fmt.Errorf(`ID=%s; URL already exists`, ID)
	}

	s.URLsMap[ID] = URL

	return nil
}

func (s Storage) Get(ID string) (string, error) {
	URL := s.URLsMap[ID]

	if URL == "" {
		return "", fmt.Errorf("URL not found")
	}

	return URL, nil
}

func GetStorage() *Storage {
	return &Storage{make(map[string]string)}
}
