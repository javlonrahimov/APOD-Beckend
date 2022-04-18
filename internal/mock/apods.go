package mock

import (
	"javlonrahimov/apod/internal/data"
)

type apodModelMock struct{}

func NewApodsMock() data.ApodService {
	return &apodModelMock{}
}

func (a *apodModelMock) Insert(apod *data.Apod) error {
	// query := `
	// INSERT INTO apods (title, explanation, date, media_type, url, hd_url)
	// VALUES ($1, $2, $3, $4, $5, $6)
	// RETURNING id, version`

	// args := []interface{}{apod.Title, apod.Explanation, apod.Date, apod.MediaType, apod.Url, apod.HdUrl}

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// return a.db.QueryRowContext(ctx, query, args...).Scan(&apod.ID, &apod.Version)
	return nil
}

func (a *apodModelMock) GetById(id int64) (*data.Apod, error) {

	if id > 0 {
		return nil, data.ErrRecordNotFound
	}

	// query := `
	// SELECT id, title, explanation, date, media_type, url, hd_url, version
	// FROM apods WHERE id = $1`

	// var apod Apod

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// err := a.db.QueryRowContext(ctx, query, id).Scan(
	// 	&apod.ID,
	// 	&apod.Title,
	// 	&apod.Explanation,
	// 	&apod.Date,
	// 	&apod.MediaType,
	// 	&apod.Url,
	// 	&apod.HdUrl,
	// 	&apod.Version,
	// )

	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return nil, data.ErrRecordNotFound
	// 	default:
	// 		return nil, err
	// 	}
	// }

	// return &apod, nil
	return nil, nil
}

func (a *apodModelMock) GetByDate(date string) (*data.Apod, error) {

	if date == "" {
		return nil, data.ErrRecordNotFound
	}

	// query := `
	// SELECT id, title, explanation, date, media_type, url, hd_url, version
	// FROM apods WHERE id = $1`

	// var apod Apod

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// err := a.db.QueryRowContext(ctx, query, date).Scan(
	// 	&apod.ID,
	// 	&apod.Title,
	// 	&apod.Explanation,
	// 	&apod.Date,
	// 	&apod.MediaType,
	// 	&apod.Url,
	// 	&apod.HdUrl,
	// 	&apod.Version,
	// )

	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return nil, ErrRecordNotFound
	// 	default:
	// 		return nil, err
	// 	}
	// }

	// return &apod, nil

	return nil, nil
}

func (a *apodModelMock) Update(apod *data.Apod) error {

	// query := `
	// UPDATE apods
	// SET title = $1, explanation = $2, date = $3, media_type = $4, url = $6, hd_url = $7, version = version + 1
	// WHERE id = $8 AND version = $9
	// RETURNING version`

	// args := []interface{}{
	// 	apod.Title,
	// 	apod.Explanation,
	// 	apod.Date,
	// 	apod.MediaType,
	// 	apod.Url,
	// 	apod.HdUrl,
	// 	apod.ID,
	// 	apod.Version,
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// err := a.db.QueryRowContext(ctx, query, args...).Scan(&apod.Version)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return ErrEditConflict
	// 	default:
	// 		return err
	// 	}
	// }
	return nil
}

func (m *apodModelMock) Delete(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	// query := `
	// DELETE FROM apods
	// WHERE id = $1`

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// result, err := m.db.ExecContext(ctx, query, id)
	// if err != nil {
	// 	return err
	// }

	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	// if rowsAffected == 0 {
	// 	return ErrRecordNotFound
	// }

	return nil
}

func (a *apodModelMock) GetAll(title string, filters data.Filters) ([]*data.Apod, data.Metadata, error) {
	// query := fmt.Sprintf(`
	// 	SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
	// 	FROM movies
	// 	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
	// 	ORDER BY %s %s, id ASC
	// 	LIMIT $2 OFFSET $3`, filters.SortColumn(), filters.SortDirection())

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	// args := []interface{}{title, filters.Limit(), filters.Offset()}

	// rows, err := a.db.QueryContext(ctx, query, args...)
	// if err != nil {
	// 	return nil, Metadata{}, err
	// }

	// defer rows.Close()

	// totalRecords := 0
	// apods := []*Apod{}

	// for rows.Next() {
	// 	var apod Apod

	// 	err := rows.Scan(
	// 		&totalRecords,
	// 		&apod.ID,
	// 		&apod.Title,
	// 		&apod.Explanation,
	// 		&apod.MediaType,
	// 		&apod.Date,
	// 		&apod.Url,
	// 		&apod.HdUrl,
	// 		&apod.Version,
	// 	)

	// 	if err != nil {
	// 		return nil, Metadata{}, err
	// 	}

	// 	apods = append(apods, &apod)
	// }

	// metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// return apods, metadata, nil

	return nil, data.Metadata{}, nil
}

//TODO: write validation
