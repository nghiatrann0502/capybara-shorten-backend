package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/nghiatrann0502/capybara-shorten-backend/internal/url-shorten/model"
	"time"
)

const dbTimeout = 3 * time.Second

type Store struct {
	db *sql.DB
}

func New() (*Store, error) {
	db, err := sql.Open("mysql", "url_shorten_dev:my_secret@tcp(localhost:3306)/url_shorten")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	store := &Store{
		db,
	}

	return store, nil
}

func (store *Store) CreateShortURL(data model.CreateShorten) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO url (short_id, original_url, created_at, updated_at)
		values (?, ?, ?, ?)`

	insertResult, err := store.db.ExecContext(ctx, stmt,
		data.ShortId,
		data.Url,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (store *Store) GetLongUrl(shortId string) (string, error) {
	var originalUrl string
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT original_url FROM url WHERE short_id = ?`

	err := store.db.QueryRowContext(ctx, stmt, shortId).Scan(&originalUrl)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("URL not found")
		}

		return "", err
	}

	return originalUrl, nil
}

func (store *Store) Close() error {
	return store.db.Close()
}
