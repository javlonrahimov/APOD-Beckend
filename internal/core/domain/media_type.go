package domain

import (
	"errors"
	"strconv"
)

var ErrInvaliMediaTypeFormat = errors.New("invalid media_type format")

type MediaType int

const (
	Image MediaType = iota
	Video
)

func (mt MediaType) MarshalJSON() ([]byte, error) {

	var jsonValue string = "unknown"

	switch mt {
	case Image:
		jsonValue = "image"
	case Video:
		jsonValue = "video"
	}

	return []byte(strconv.Quote(jsonValue)), nil
}

func (mt *MediaType) UnmarshalJSON(jsonValue []byte) error {

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvaliMediaTypeFormat
	}

	switch unquotedJSONValue {
	case "image":
		*mt = Image
	case "video":
		*mt = Video
	default:
		return ErrInvaliMediaTypeFormat
	}
	return nil
}
