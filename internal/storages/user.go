package storages

import "github.com/jackc/pgx"

type UserStorage interface {
	InsertUser()
	GetFullUser()
	UpdateUser()

	GetAllUsersByForum()
}

type userStorage struct {
	db *pgx.ConnPool
}

/*func NewUserStorage(db *pgx.ConnPool) UserStorage {
	return &userStorage{
		db: db,
	}
}*/