package data

import (
	"database/sql"
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

func (a ApodModel) Insert(movie *Apod) error {
	return nil
}

func (a ApodModel) Get(id int64) (*Apod, error) {
	return nil, nil
}

func (a ApodModel) Update(movie *Apod) error {
	return nil
}

func (a ApodModel) Delete(id int64) error {
	return nil
}

func ValidateApod(v *validator.Validator, apod *Apod) {
	v.Check(apod.Title != "", "title", "must be provided")
	v.Check(len(apod.Title) <= 500, "title", "must not be longer than 500 bytes long")
}
