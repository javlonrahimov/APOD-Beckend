package data

import (
	"context"
	"database/sql"
	"time"
)

type LikesService interface {
	Insert(apodId, userId int64) error
	Delete(apodId, userId int64) error
}

type Like struct {
	Id     int64
	UserId int64
	ApodId int64
}

type likeModel struct {
	db *sql.DB
}

func NewLikeModel(db *sql.DB) LikesService {
	return &likeModel{db: db}
}

func (l likeModel) Delete(apodId int64, userId int64) error {
	query := `
	DELETE FROM likes
	WHERE apod_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := l.db.ExecContext(ctx, query, apodId, userId)

	return err
}

func (l likeModel) Insert(apodId int64, userId int64) error {

	query := `
	INSERT INTO likes (user_id, apod_id)
	VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := l.db.ExecContext(ctx, query, userId, apodId)

	return err
}
