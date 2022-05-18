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

type DatabaseStorage struct {
	DB *pgx.Conn
}

type Repository interface {
	AddURL(UserShortURL) error
	GetURL(string) (string, error)
	GetUserShortURLs(string) []UserShortURL
	Destruct(app.Config) error
	AddBatchURL([]UserShortURL) error
	GetShortByOriginal(string) (string, error)
	Ping() error
}

func New(cfg app.Config) Repository {
	if cfg.DatabaseDSN != "" {
		dbStore, err  := constructDatabaseStorage(cfg)
		if err != nil {
			fmt.Println("constructLocalStorage ERROR: ", err)
			return constructLocalStorage(cfg)
		}

		return dbStore
	}

	return constructLocalStorage(cfg)
}
