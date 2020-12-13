package storages

import "github.com/jackc/pgx"

type PostStorage interface {
	InsertSomePosts()
	InsertPost()
	UpdatePost()
	GetFullPost()

	GetAllPostsByThread()
}

type postStorage struct {
	db *pgx.ConnPool
}

/*func NewPostStorage(db *pgx.ConnPool) PostStorage{
	return &postStorage{
		db: db,
	}
}*/