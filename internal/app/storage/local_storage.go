package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/Anav11/url-shortener/internal/app"
)

type LocalStorage struct {
	URLsMap URLsMap
	UserURLs UserURLs
	mutex   sync.RWMutex
}

func (ls *LocalStorage) AddURL(usu UserShortURL) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	if ls.URLsMap[usu.ID] != "" {
		return fmt.Errorf(`ID=%s; URL already exists`, usu.ID)
	}

	existID, _ := ls.GetShortByOriginal(usu.OriginalURL)
	if existID != "" {
		return &URLDuplicateError{URL: usu.OriginalURL}
	}

	ls.URLsMap[usu.ID] = usu.OriginalURL
	ls.UserURLs[usu.UserID] = append(ls.UserURLs[usu.UserID], usu.ID)

	return nil
}

func (ls *LocalStorage) GetURL(ID string) (string, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	URL := ls.URLsMap[ID]
	if URL == "" {
		return "", fmt.Errorf("URL not found")
	}

	return URL, nil
}

func (ls *LocalStorage) GetUserShortURLs(userID string) []UserShortURL {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	userURLs := ls.UserURLs[userID]

	var URLs []UserShortURL
	for _, ID := range userURLs {
		URL, _ := ls.GetURL(ID)
		URLs = append(URLs, UserShortURL{ID, URL, userID})
	}

	return URLs
}

func (ls *LocalStorage) AddBatchURL(URLs []UserShortURL) error {
	for _, URL := range URLs {
		err := ls.AddURL(URL)
		if err != nil {
				return err
			}
	}

	return nil
}

func (ls *LocalStorage) GetShortByOriginal(originalURL string) (string, error) {
	for ID, URL := range ls.URLsMap {
		if URL == originalURL {
			return ID, nil
		}
	}

	return "", fmt.Errorf("URL not found")
}

func (ls *LocalStorage) DeleteUserURLs(IDs []string, userID string) error {
	return nil
}

func constructLocalStorage(cfg app.Config) Repository {
	ls := &LocalStorage{make(URLsMap), make(UserURLs), sync.RWMutex{}}

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("OpenFile error; %s", err)
		return ls
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Printf("ReadAll error; %s", err)
		return ls
	}

	for _, rec := range records {
		ls.URLsMap[rec[0]] = rec[1]
	}

	return ls
}

func (ls *LocalStorage) Destruct(cfg app.Config) error {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_WRONLY, 0664)
	if err != nil {
		return fmt.Errorf("OpenFile error; %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	var records [][]string
	for ID, URL := range ls.URLsMap {
		records = append(records, []string{ID, URL})
	}

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("WriteAll error; %s", err)
	}

	writer.Flush()

	return nil
}

func (ls *LocalStorage) Ping() error {
	return nil
}