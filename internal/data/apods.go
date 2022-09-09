package data

import (
	"context"
	"database/sql"
	"errors"
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
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Version     int       `json:"-"`
}

type ApodModel struct {
	DB *sql.DB
}

func (a ApodModel) Insert(apod *Apod) error {

	query := `INSERT INTO apods (date, explanation, hd_url, url, title, media_type)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at`

	args := []interface{}{apod.Date, apod.Explanation, apod.HdUrl, apod.Url, apod.Title, apod.MediaType}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&apod.ID, &apod.CreatedAt)
}

func (a ApodModel) Get(id int64) (*Apod, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, date, explanation, hd_url, url, title, media_type
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
	SET date = $1, explanation = $2, hd_url = $3, url = $4, title = $5, media_type = $6, version = version + 1, updated_at = CURRENT_TIMESTAMP
	WHERE id = $7 AND version = $8
	RETURNING version`

	args := []interface{}{
		apod.Date, apod.Explanation, apod.HdUrl, apod.Url, apod.Title, apod.MediaType, apod.ID, apod.Version,
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

func ValidateApod(v *validator.Validator, apod *Apod) {
	v.Check(apod.Title != "", "title", "must be provided")
	v.Check(len(apod.Title) <= 500, "title", "must not be longer than 500 bytes long")
}
