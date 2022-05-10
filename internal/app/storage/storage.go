package storage

import (
	"fmt"
	"github.com/Anav11/url-shortener/internal/app"
	"github.com/jackc/pgx/v4"
)

type URLsMap = map[string]string
type UserURLs = map[string][]string

type UserShortURL struct {
	ID 			string
	OriginalURL string
	UserID      string
}

type Repository interface {
	AddURL(UserShortURL) error
	GetURL(string) (string, error)
	GetUserShortURLs(string) []UserShortURL
	DestructStorage(app.Config) error
	AddBatchURL([]UserShortURL) error
	GetShortByOriginal(string) (string, error)
	Ping() error
}

type DatabaseStorage struct {
	DB *pgx.Conn
}

func ConstructStorage(conf app.Config) Repository {
	if conf.DatabaseDSN != "" {
		dbStore, err  := ConstructDatabaseStorage(conf)
		if err != nil {
			fmt.Println("ConstructDatabaseStorage ERROR: ", err)
			return ConstructLocalStorage(conf)
		}

		return dbStore
	}

	return ConstructLocalStorage(conf)
}
