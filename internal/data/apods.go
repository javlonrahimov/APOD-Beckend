package data

import "time"

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
