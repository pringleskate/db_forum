package storages

import (
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
)

type PostStorage interface {
	InsertSomePosts(inputPosts []models.Post) (posts []models.Post, err error)
	InsertPost(postInput models.Post) (ID int, err error)
	UpdatePost(postInput models.PostUpdate) (err error)
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

func (p postStorage) InsertSomePosts(inputPosts []models.Post) (posts []models.Post, err error) {
	for i, post := range inputPosts {
		ID, err := p.InsertPost(post)
		if err != nil {
			return
		}

		inputPosts[i].ID = ID
	}

	return inputPosts, nil
}

func (p postStorage) InsertPost(postInput models.Post) (ID int, err error) {
	if postInput.Parent == 0 {
		err = p.db.QueryRow("INSERT INTO post (author, created, forum, message, parent, thread, path) VALUES ($1,$2,$3,$4,$5,$6, array[(select currval('post_id_seq')::integer)]) RETURNING ID",
			postInput.Author, postInput.Created, postInput.Forum, postInput.Message, postInput.Parent, postInput.ThreadSlagOrID.ThreadID).Scan(&ID)

	} else {
		err = p.db.QueryRow("INSERT INTO post (author, created, forum, message, parent, thread, path) VALUES ($1,$2,$3,$4,$5,$6, (SELECT path FROM post WHERE id = $5) || (select currval('post_id_seq')::integer)) RETURNING ID",
			postInput.Author, postInput.Created, postInput.Forum, postInput.Message, postInput.Parent, postInput.ThreadSlagOrID.ThreadID).Scan(&ID)
	}
	return
}

func (p postStorage) UpdatePost(postInput models.PostUpdate) (err error) {

}

func (p postStorage) GetFullPost() {
	panic("implement me")
}

func (p postStorage) GetAllPostsByThread() {
	panic("implement me")
}
