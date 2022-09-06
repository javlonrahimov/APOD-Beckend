package data

import (
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
}

type ApodModel struct {
	DB *sql.DB
}

func (a ApodModel) Insert(apod *Apod) error {

	query := `INSERT INTO apods (date, explanation, hd_url, url, title, media_type)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at`

	args := []interface{}{apod.Date, apod.Explanation, apod.HdUrl, apod.Url, apod.Title, apod.MediaType}

	return a.DB.QueryRow(query, args...).Scan(&apod.ID, &apod.CreatedAt)
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

	err := a.DB.QueryRow(query, id).Scan(
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
		default :
			return nil,err
		}
	}

	return &apod, nil
}

func (a ApodModel) Update(apod *Apod) error {
	return nil
}

func (a ApodModel) Delete(id int64) error {
	return nil
}

func ValidateApod(v *validator.Validator, apod *Apod) {
	v.Check(apod.Title != "", "title", "must be provided")
	v.Check(len(apod.Title) <= 500, "title", "must not be longer than 500 bytes long")
}
