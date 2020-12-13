package storages

import "github.com/jackc/pgx"

type ForumStorage interface {
	InsertForum()
	GetFullForum()
	InsertForumUser()
	UpdateThreadsCount()
	UpdatePostsCount()
	Status()
	Clear()
}

type forumStorage struct {
	db *pgx.ConnPool
}

/*func NewForumStorage(db *pgx.ConnPool) ForumStorage {
	return &forumStorage{
		db: db,
	}
}*/