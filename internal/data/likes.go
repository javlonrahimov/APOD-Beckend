package data

import "database/sql"

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

// Delete implements LikesService
func (*likeModel) Delete(apodId int64, userId int64) error {
	panic("unimplemented")
}

// Insert implements LikesService
func (*likeModel) Insert(apodId int64, userId int64) error {
	panic("unimplemented")
}
