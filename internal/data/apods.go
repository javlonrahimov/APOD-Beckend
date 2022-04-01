package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Apod struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Date        string `json:"date"`
	MediaType   string `json:"media_type"`
	Url         string `json:"url"`
	HdUrl       string `json:"hd_url"`
	Version     int32  `json:"version"`
}

type ApodModel struct {
	DB *sql.DB
}

func (m ApodModel) Insert(apod * Apod) error {
	query := `
	INSERT INTO apods (title, explanation, date, media_type, url, hd_url)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version`

	args := []interface{}{apod.Title, apod.Explanation, apod.Date, apod.MediaType, apod.Url, apod.HdUrl}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx ,query, args...).Scan(&apod.ID, &apod.Version)
}


func (m ApodModel) GetById(id int64) (*Apod, error) {

	if id > 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, title, explanation, date, media_type, url, hd_url, version
	FROM apods WHERE id = $1`

	var apod Apod

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&apod.ID,
		&apod.Title,
		&apod.Explanation,
		&apod.Date,
		&apod.MediaType,
		&apod.Url,
		&apod.HdUrl,
		&apod.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &apod, nil
}

func (m ApodModel) GetByDate(date string) (*Apod, error) {

	if date == "" {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, title, explanation, date, media_type, url, hd_url, version
	FROM apods WHERE id = $1`

	var apod Apod

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, date).Scan(
		&apod.ID,
		&apod.Title,
		&apod.Explanation,
		&apod.Date,
		&apod.MediaType,
		&apod.Url,
		&apod.HdUrl,
		&apod.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &apod, nil
}

func (m ApodModel) Update(apod * Apod) error {

	query := `
	UPDATE apods
	SET title = $1, explanation = $2, date = $3, media_type = $4, url = $6, hd_url = $7, version = version + 1
	WHERE id = $8 AND version = $9
	RETURNING version`

	args := []interface{}{
		apod.Title,
		apod.Explanation,
		apod.Date,
		apod.MediaType,
		apod.Url,
		apod.HdUrl,
		apod.ID,
		apod.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&apod.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m ApodModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM apods
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
