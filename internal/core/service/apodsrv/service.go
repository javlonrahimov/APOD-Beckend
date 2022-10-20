package apodsrv

import (
	"time"

	"apod.api.javlonrahimov1212/internal/core/domain"
	"apod.api.javlonrahimov1212/internal/core/ports"
	"apod.api.javlonrahimov1212/internal/data"
	"apod.api.javlonrahimov1212/internal/validator"
)

type service struct {
	apodRepository ports.ApodRepository
}

func New(apodRepository ports.ApodRepository) *service {
	return &service{
		apodRepository: apodRepository,
	}
}

func (srv *service) Create(title, date, explanation, hdUrl, url, copyright, mediaType string, validator *validator.Validator ) (*domain.Apod, error) {
	apod := &domain.Apod{
		Date:        time.Time{},
		Explanation: explanation,
		HdUrl:       hdUrl,
		Url:         url,
		Title:       title,
		MediaType:   domain.GetMediaType(mediaType),
		Copyright:   copyright,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	
	if domain.ValidateApod(validator, apod); !validator.Valid() {
		return nil, nil
	}

	err := srv.apodRepository.Insert(apod)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (srv *service) Get(id int64) (domain.Apod, error) {
	apod, err := srv.apodRepository.Get(id)
	return *apod, err
}

func (srv *service) Update(id, title, date, explanation, hdUrl, url, mediaType string) (domain.Apod, error) {

}

func (srv *service) Delete(id string) error {

}

func (srv *service) GetPaging(title, copyright string, date time.Time, filters data.Filters) {

}
