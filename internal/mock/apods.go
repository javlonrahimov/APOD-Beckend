package mock

import (
	"javlonrahimov/apod/internal/data"
	"strings"
)

type apodModelMock struct{}

var apods = make([]data.Apod, 0)

func NewApodsMock() data.ApodService {
	return &apodModelMock{}
}

func (a *apodModelMock) Insert(apod *data.Apod) error {
	apods = append(apods, *apod)
	return nil
}

func (a *apodModelMock) GetById(id int64) (*data.Apod, error) {

	if id > 0 {
		return nil, data.ErrRecordNotFound
	}

	for i := 0; i < len(apods); i++ {
		if apods[i].ID == id {
			return &apods[i], nil
		}
	}

	return nil, data.ErrRecordNotFound
}

func (a *apodModelMock) GetByDate(date string) (*data.Apod, error) {

	if date == "" {
		return nil, data.ErrRecordNotFound
	}

	for i := 0; i < len(apods); i++ {
		if apods[i].Date == date {
			return &apods[i], nil
		}
	}

	return nil, data.ErrRecordNotFound
}

func (a *apodModelMock) Update(apod *data.Apod) error {

	for i := 0; i < len(apods); i++ {
		if apods[i].ID == apod.ID {
			apods[i] = *apod
			return nil
		}
	}

	return nil
}

func (m *apodModelMock) Delete(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	for i := 0; i < len(apods); i++ {
		if apods[i].ID == id {
			apods = removeApod(apods, i)
			return nil
		}
	}

	return data.ErrRecordNotFound
}

func (a *apodModelMock) GetAll(title string, filters data.Filters) ([]*data.Apod, data.Metadata, error) {

	pagingApods := make([]*data.Apod, 0)

	for i := 0; i < len(apods); i++ {
		if strings.Contains(apods[i].Title, title) {
			pagingApods = append(pagingApods, &apods[i])
		}
	}

	return pagingApods, data.Metadata{}, nil
}

func removeApod(slice []data.Apod, s int) []data.Apod {
	return append(slice[:s], slice[s+1:]...)
}