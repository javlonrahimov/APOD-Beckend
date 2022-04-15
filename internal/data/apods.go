package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ApodService interface {
		Insert(apod *Apod) error
		GetById(id int64) (*Apod, error)
		GetByDate(date string) (*Apod, error)
		Update(apod *Apod) error
		Delete(id int64) error
		GetAll(title string, filters Filters) ([]*Apod, Metadata, error)
}

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

type apodModel struct {
	db *sql.DB
}

func NewApodModel(db *sql.DB) ApodService {
	return &apodModel{db: db}
}


func (a *apodModel) Insert(apod *Apod) error {
	query := `
	INSERT INTO apods (title, explanation, date, media_type, url, hd_url)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version`

	args := []interface{}{apod.Title, apod.Explanation, apod.Date, apod.MediaType, apod.Url, apod.HdUrl}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.db.QueryRowContext(ctx, query, args...).Scan(&apod.ID, &apod.Version)
}

func (a *apodModel) GetById(id int64) (*Apod, error) {

	if id > 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, title, explanation, date, media_type, url, hd_url, version
	FROM apods WHERE id = $1`

	var apod Apod

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.db.QueryRowContext(ctx, query, id).Scan(
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

func (a *apodModel) GetByDate(date string) (*Apod, error) {

	if date == "" {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, title, explanation, date, media_type, url, hd_url, version
	FROM apods WHERE id = $1`

	var apod Apod

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.db.QueryRowContext(ctx, query, date).Scan(
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

func (a *apodModel) Update(apod *Apod) error {

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.db.QueryRowContext(ctx, query, args...).Scan(&apod.Version)
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

func (m *apodModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM apods
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.db.ExecContext(ctx, query, id)
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

func (a *apodModel) GetAll(title string, filters Filters) ([]*Apod, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, filters.limit(), filters.offset()}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	apods := []*Apod{}

	for rows.Next() {
		var apod Apod

		err := rows.Scan(
			&totalRecords,
			&apod.ID,
			&apod.Title,
			&apod.Explanation,
			&apod.MediaType,
			&apod.Date,
			&apod.Url,
			&apod.HdUrl,
			&apod.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		apods = append(apods, &apod)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return apods, metadata, nil
}

//TODO: write validation
