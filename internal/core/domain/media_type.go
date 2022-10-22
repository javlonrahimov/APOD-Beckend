package domain

import (
	"strconv"
)

type MediaType int

const (
	Image MediaType = iota
	Video
	Unknown
)

func GetMediaType(data string) MediaType {
	switch data {
	case "image":
		return Image
	case "video":
		return Video
	default:
		return Unknown
	}
}

func (mt MediaType) MarshalJSON() ([]byte, error) {

	var jsonValue string = "unknown"

	switch mt {
	case Image:
		jsonValue = "image"
	case Video:
		jsonValue = "video"
	case Unknown:
		jsonValue = "unknown"
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
		*mt = Unknown
	}
	return nil
}
