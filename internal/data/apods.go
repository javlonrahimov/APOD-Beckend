package data

import (
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

func ValidateApod(v *validator.Validator, apod *Apod) {
	v.Check(apod.Title != "", "title", "must be provided")
	v.Check(len(apod.Title) <= 500, "title", "must not be longer than 500 bytes long")
}
