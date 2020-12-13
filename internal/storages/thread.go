package storages

import "github.com/jackc/pgx"

type ThreadStorage interface {
	InsertThread()
	UpdateThread()
	GetFullThread()

	InsertVote()
	UpdateVote()
	GetAllThreadsByForum()
}

type threadStorage struct {
	db *pgx.ConnPool
}
/*

func NewThreadStorage(db *pgx.ConnPool) ThreadStorage{
	return &threadStorage{
		db: db,
	}
}*/