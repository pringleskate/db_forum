package storages

import (
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
)

type PostStorage interface {
	InsertSomePosts()
	InsertPost(postInput models.Post) (err error)
	UpdatePost()
	GetFullPost()
	GetAllPostsByThread()
}

type postStorage struct {
	db *pgx.ConnPool
}

func NewPostStorage(db *pgx.ConnPool) PostStorage{
	return &postStorage{
		db: db,
	}
}

func (p postStorage) InsertSomePosts() {
	panic("implement me")
}

func (p postStorage) InsertPost(postInput models.Post) (err error) {
	return
}

func (p postStorage) UpdatePost() {
	panic("implement me")
}

func (p postStorage) GetFullPost() {
	panic("implement me")
}

func (p postStorage) GetAllPostsByThread() {
	panic("implement me")
}
