package apodsrv

import (
	"apod.api.javlonrahimov1212/internal/core/domain"
	"apod.api.javlonrahimov1212/internal/core/ports"
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

func (srv *service) Create(title, date, explanation, hdUrl, url, copyright, mediaType string) (*domain.Apod, *validator.Validator, error) {

	validator := validator.New()

	apod := &domain.Apod{
		Date:        date,
		Explanation: explanation,
		HdUrl:       hdUrl,
		Url:         url,
		Title:       title,
		MediaType:   domain.GetMediaType(mediaType),
		Copyright:   copyright,
	}
	
	if domain.ValidateApod(validator, apod); !validator.Valid() {
		return nil, validator, nil
	}

	err := srv.apodRepository.Insert(apod)
	if err != nil {
		return nil, nil, err
	}
	return apod, nil, nil
}

func (srv *service) Get(id int64) (*domain.Apod, error) {
	return srv.apodRepository.Get(id)
}

func (srv *service) Update(id int64, title, date, explanation, hdUrl, url, mediaType string) (*domain.Apod, *validator.Validator, error) {
	apod, err := srv.apodRepository.Get(id)
	if err != nil {
		return nil, nil, err
	}

	input := struct {
		Date        *string
		Explanation *string
		HdUrl       *string
		Url         *string
		Title       *string
		MediaType   domain.MediaType
	}{
		Date:        &date,
		Explanation: &explanation,
		HdUrl:       &hdUrl,
		Url:         &url,
		Title:       &title,
		MediaType:   domain.GetMediaType(mediaType),
	}


	if input.Date != nil {
		apod.Date = *input.Date
	}

	if input.Explanation != nil {
		apod.Explanation = *input.Explanation
	}

	if input.HdUrl != nil {
		apod.HdUrl = *input.HdUrl
	}

	if input.Url != nil {
		apod.Url = *input.Url
	}

	if input.Title != nil {
		apod.Title = *input.Title
	}

	if input.MediaType != domain.Unknown {
		apod.MediaType = input.MediaType
	}

	v := validator.New()

	if domain.ValidateApod(v, apod); !v.Valid() {
		return nil, v, nil
	}

	err = srv.apodRepository.Update(apod)
	if err != nil {
		return nil, nil, err
	}

	return apod, nil, nil
}

func (srv *service) Delete(id int64) error {
	err := srv.apodRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (srv *service) GetPaging(title, copyright, date string, filters domain.Filters, validator *validator.Validator) ([]*domain.Apod, domain.Metadata, error) {
	
	if domain.ValidateFilters(validator, filters); !validator.Valid() {
		return nil, domain.Metadata{}, nil
	}

	apods, metadata, err := srv.apodRepository.GetAll(title, copyright, date, filters)
	if err != nil {
		return nil, domain.Metadata{}, err
	}

	return apods, metadata, nil
}
