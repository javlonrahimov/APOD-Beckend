package domain

import (
	"time"

	"apod.api.javlonrahimov1212/internal/validator"
)

type Apod struct {
	ID          int64     `json:"id"`
	Date        string 	  `json:"date,omitempty"`
	Explanation string    `json:"explanation,omitempty"`
	HdUrl       string    `json:"hd_url,omitempty"`
	Url         string    `json:"url,omitempty"`
	Title       string    `json:"title,omitempty"`
	MediaType   MediaType `json:"media_type"`
	Copyright   string    `json:"copyright"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Version     int       `json:"-"`
}

func ValidateApod(v *validator.Validator, apod *Apod) {
	v.Check(apod.MediaType != Unknown, "media_type", "invalid media type")
	v.Check(apod.Title != "", "title", "must be provided")
	v.Check(len(apod.Title) <= 500, "title", "must not be longer than 500 bytes long")
}
