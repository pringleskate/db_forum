package storages

import (
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
)

type UserStorage interface {
	InsertUser(userInput models.User) (err error)
	GetFullUser(nickname string) (user models.User, err error)
	UpdateUser(userInput models.User) (err error)
	GetAllUsersByForum(params models.ForumQueryParams) (users []models.User, err error)
}

type userStorage struct {
	db *pgx.ConnPool
}

func NewUserStorage(db *pgx.ConnPool) UserStorage {
	return &userStorage{
		db: db,
	}
}

func (u userStorage) InsertUser(userInput models.User) (err error){
	_, err = u.db.Exec("INSERT INTO users (nick_name, email, full_name, about) VALUES ($1, $2, $3, $4)",
		userInput.Nickname, userInput.Email, userInput.FullName, userInput.About)
	return
}

func (u userStorage) GetFullUser(nickname string) (user models.User, err error) {
	err = u.db.QueryRow("SELECT nick_name, email, full_name, about FROM users WHERE lower(nick_name) = lower($1)", nickname).
		Scan(&user.Nickname, &user.Email, &user.FullName, &user.About)
	return
}

func (u userStorage) UpdateUser(userInput models.User) (err error) {
	_, err = u.db.Exec("UPDATE users SET email = $1, full_name = $2, about = $3 WHERE lower(nick_name) = lower($4)",
		userInput.Email, userInput.FullName, userInput.About, userInput.Nickname)
	return
}

func (u userStorage) GetAllUsersByForum(params models.ForumQueryParams) (users []models.User, err error) {
	var rows *pgx.Rows

	if params.Desc {
		if params.Since == "" {
			//	selectWithDesc = "SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users fu JOIN users u ON fu.userID = u.ID WHERE fu.forumID = $1 ORDER BY u.nickname DESC LIMIT $2"
			rows, err = u.db.Query(`SELECT u.nick_name, u.email, u.full_name, u.about FROM users u
				JOIN forum_users fu ON fu.user = u.nick_name
				WHERE lower(fu.forum) = lower($1) ORDER BY lower(u.nick_name) COLLATE "C" DESC LIMIT $2`,
				params.Slug, params.Limit)
		} else {
			rows, err =  u.db.Query(`SELECT u.nick_name, u.email, u.full_name, u.about FROM users u 
				JOIN forum_users fu ON fu.user = u.nick_name
				WHERE lower(fu.forum) = lower($1) AND lower(u.nick_name) < lower($3) ORDER BY lower(u.nick_name) COLLATE "C" DESC LIMIT $2`,
				params.Slug, params.Limit, params.Since)
		}
	} else {
		if params.Since == "" {
			rows, err = u.db.Query(`SELECT u.nick_name, u.email, u.full_name, u.about FROM users u
				JOIN forum_users fu ON fu.user = u.nick_name
				WHERE lower(fu.forum) = lower($1) ORDER BY lower(u.nick_name) COLLATE "C" LIMIT $2`,
				params.Slug, params.Limit)
		} else {
			rows, err =  u.db.Query(`SELECT u.nick_name, u.email, u.full_name, u.about FROM users u 
				JOIN forum_users fu ON fu.user = u.nick_name
				WHERE lower(fu.forum) = lower($1) AND lower(u.nick_name) < lower($3) ORDER BY lower(u.nick_name) COLLATE "C" LIMIT $2`,
				params.Slug, params.Limit, params.Since)
		}
	}

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}

		err = rows.Scan(&user.Nickname, &user.Email, &user.FullName, &user.About)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}
