package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/Anav11/url-shortener/internal/app"
)

type URLDuplicateError struct {
	URL string
}

func (err *URLDuplicateError) Error() string {
	return fmt.Sprintf("URL %s - already exists.", err.URL)
}

func (dbs *DatabaseStorage) AddURL(usu UserShortURL) error {
	_, err := dbs.DB.Exec(context.Background(), "INSERT INTO urls VALUES ($1, $2, $3)", usu.ID, usu.OriginalURL, usu.UserID)

	var pgError *pgconn.PgError

	if errors.As(err, &pgError) {
		if pgError.Code == pgerrcode.UniqueViolation {
			return &URLDuplicateError{URL: usu.OriginalURL}
		}
	}

	return err
}

func (dbs *DatabaseStorage) GetURL(ID string) (string, error) {
	row := ""
	err := dbs.DB.QueryRow(context.Background(), "SELECT original_url FROM urls WHERE url_id = $1 AND is_deleted = false", ID).Scan(&row)
	if err != nil {
		return "", err
	}

	return row, nil
}

func (dbs *DatabaseStorage) GetUserShortURLs(userID string) []UserShortURL {
	shortURLs := make([]UserShortURL, 0)
	rows, err := dbs.DB.Query(context.Background(), "SELECT url_id, original_url FROM urls WHERE user_id = $1 AND is_deleted = false", userID)
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

func (dbs *DatabaseStorage) AddBatchURL(shortURLs []UserShortURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stmt, err := dbs.DB.Prepare(ctx, "addBatch", "INSERT INTO urls VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	for _, su := range shortURLs {
		_, err := dbs.DB.Exec(ctx, stmt.SQL, su.ID, su.OriginalURL, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (dbs *DatabaseStorage) GetShortByOriginal(originalURL string) (string, error) {
	var ID string
	if err := dbs.DB.QueryRow(context.Background(), "SELECT url_id FROM urls WHERE original_url = $1", originalURL).Scan(&ID); err != nil {
		return "", err
	}

	return ID, nil
}

func (dbs *DatabaseStorage) DeleteUserURLs(IDs []string, userID string) error {
	_, err := dbs.DB.Exec(context.Background(), "UPDATE urls SET is_deleted = true WHERE user_id = $1 AND url_id = ANY($2)", userID, IDs)
	return err
}

func constructDatabaseStorage(cfg app.Config) (Repository, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	dbs := &DatabaseStorage{ DB: conn }

	const createTable = `
		CREATE TABLE IF NOT EXISTS urls (
			url_id varchar(36) NOT NULL UNIQUE PRIMARY KEY,
			original_url varchar(255) UNIQUE,
			user_id varchar(36),
			is_deleted boolean DEFAULT false
		)`
	if _, err = dbs.DB.Exec(context.Background(), createTable); err != nil {
		return nil, err
	}

	return dbs, nil
}

func (dbs DatabaseStorage) Destruct(cfg app.Config) error {
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