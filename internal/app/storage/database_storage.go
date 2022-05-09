package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/Anav11/url-shortener/internal/app"
)

func (dbs *DatabaseStorage) AddURL(ID string, URL string, userID string) error {
	_, err := dbs.DB.Exec(context.Background(), "INSERT INTO urls VALUES ($1, $2, $3)", ID, URL, userID)

	return err
}

func (dbs *DatabaseStorage) GetURL(ID string) (string, error) {
	row := ""
	err := dbs.DB.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE url_id = $1", ID).Scan(&row)
	if err != nil {
		return "", err
	}

	return row, nil
}

func (dbs *DatabaseStorage) GetUserShortURLs(userID string) []UserShortURL {
	shortURLs := make([]UserShortURL, 0)
	rows, err := dbs.DB.Query(context.Background(), "SELECT url_id, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return shortURLs
	}
	defer rows.Close()

	for rows.Next() {
		var sURL UserShortURL
		err := rows.Scan(&sURL.ID, &sURL.OriginalURL)
		if err != nil {
			return nil
		}
		shortURLs = append(shortURLs, sURL)
	}

	return shortURLs
}

func ConstructDatabaseStorage(conf app.Config) (Repository, error) {
	conn, err := pgx.Connect(context.Background(), conf.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	dbs := &DatabaseStorage{ DB: conn }

	const CreateTable = `
		CREATE TABLE IF NOT EXISTS urls (
			url_id varchar(36) NOT NULL UNIQUE PRIMARY KEY,
			original_url varchar(255),
			user_id varchar(36)
		)`
	_, err = dbs.DB.Exec(context.Background(), CreateTable)
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

func (dbs DatabaseStorage) DestructStorage(conf app.Config) error {
	err := dbs.DB.Close(context.Background())

	return err
}

func (dbs DatabaseStorage) Ping() error {
	ctx := context.Background()
	conn, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := dbs.DB.Ping(conn)

	return err
}