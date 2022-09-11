package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"apod.api.javlonrahimov1212/internal/validator"
)

type Apod struct {
	ID          int64     `json:"id"`
	Date        time.Time `json:"date,ommitempty"`
	Explanation string    `json:"explanation,ommitempty"`
	HdUrl       string    `json:"hd_url,ommitempty"`
	Url         string    `json:"url,ommitempty"`
	Title       string    `json:"title,ommitempty"`
	MediaType   MediaType `json:"media_type"`
	Copyright   string    `json:"copyright"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Version     int       `json:"-"`
}

type ApodModel struct {
	DB *sql.DB
}

func (a ApodModel) Insert(apod *Apod) error {

	query := `INSERT INTO apods (date, explanation, hd_url, url, title, media_type, copyright)
	VALUES ($1, $2, $3, $4, $5, $6, &7)
	RETURNING id, created_at`

	args := []interface{}{apod.Date, apod.Explanation, apod.HdUrl, apod.Url, apod.Title, apod.MediaType, apod.Copyright}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&apod.ID, &apod.CreatedAt)
}

func (a ApodModel) Get(id int64) (*Apod, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, date, explanation, hd_url, url, title, media_type, copyright, version
		FROM apods
		WHERE id = $1`

	var apod Apod

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, id).Scan(
		&apod.ID,
		&apod.CreatedAt,
		&apod.Date,
		&apod.Explanation,
		&apod.HdUrl,
		&apod.Url,
		&apod.Title,
		&apod.MediaType,
		&apod.Copyright,
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

func (a ApodModel) Update(apod *Apod) error {

	query := `
	UPDATE apods
	SET date = $1, explanation = $2, hd_url = $3, url = $4, title = $5, media_type = $6, copyright = &7, version = version + 1, updated_at = CURRENT_TIMESTAMP
	WHERE id = $8 AND version = $9
	RETURNING version`

	args := []interface{}{
		apod.Date, apod.Explanation, apod.HdUrl, apod.Url, apod.Title, apod.MediaType, apod.Copyright, apod.ID, apod.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(&apod.Version)

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

func (a ApodModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM apods
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := a.DB.ExecContext(ctx, query, id)
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

func (a ApodModel) GetAll(title string, copyright string, date time.Time, filters Filters) ([]*Apod, Metadata, error) {

	// todo do i really need full-text search like this
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, date, explanation, hd_url, url, title, media_type, copyright
	FROM apods
	WHERE (to_tsvector('simple', title) @@  plainto_tsquery('simple', $1) OR $1 = '')
	AND (to_tsvector('simple', copyright) @@  plainto_tsquery('simple', $2) OR $2 = '')
	ORDER BY %s %s ,id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{title, copyright, filters.limit(), filters.offset()}

	rows, err := a.DB.QueryContext(ctx, query, args...)
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
			&apod.CreatedAt,
			&apod.Date,
			&apod.Explanation,
			&apod.HdUrl,
			&apod.Url,
			&apod.Title,
			&apod.MediaType,
			&apod.Copyright,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		apods = append(apods, &apod)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return apods, metadata, nil
}

func ValidateApod(v *validator.Validator, apod *Apod) {
	v.Check(apod.Title != "", "title", "must be provided")
	v.Check(len(apod.Title) <= 500, "title", "must not be longer than 500 bytes long")
}
