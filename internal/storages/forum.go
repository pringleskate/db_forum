package storages

import (
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
)

type ForumStorage interface {
	InsertForum(forumInput models.Forum) (err error)
	GetFullForum(slug string) (forum models.Forum, err error)
	InsertForumUser(forum string, user string) (err error)
	UpdateThreadsCount(forumId int) (err error)
	UpdatePostsCount(forumId int) (err error)
	Status() (status models.Status, err error)
	Clear() (err error)
}

type forumStorage struct {
	db *pgx.ConnPool
}

func NewForumStorage(db *pgx.ConnPool) ForumStorage {
	return &forumStorage{
		db: db,
	}
}

func (f forumStorage) InsertForum(forumInput models.Forum) (err error) {
	_, err = f.db.Exec("INSERT INTO forum (slug, title, author) VALUES ($1, $2, $3)", forumInput.Slug, forumInput.Title, forumInput.User)
	return
}

func (f forumStorage) GetFullForum(slug string) (forum models.Forum, err error) {
	err = f.db.QueryRow("SELECT  slug, title, threads, posts, author FROM forum WHERE lower(slug) = lower($1)", slug).
				Scan(&forum.Slug, &forum.Title, &forum.Threads, &forum.Posts, &forum.User)
	return
}

func (f forumStorage) InsertForumUser(forum string, user string) (err error) {
	_, err = f.db.Exec("INSERT INTO forum_users (forum, user) VALUES ($1, $2)", forum, user)
	return
}

func (f forumStorage) UpdateThreadsCount(forumId int) (err error) {
	_, err = f.db.Exec("UPDATE forum SET threads = threads + 1 WHERE ID = $1", forumId)
	return
}

func (f forumStorage) UpdatePostsCount(forumId int) (err error) {
	_, err = f.db.Exec("UPDATE forum SET posts = posts + 1 WHERE ID = $1", forumId)
	return
}

func (f forumStorage) Status() (status models.Status, err error){
	err = f.db.QueryRow("SELECT * FROM " +
		"(SELECT COUNT(*) FROM forum) AS F, " +
		"(SELECT COUNT(*) FROM thread) AS T," +
		"(SELECT COUNT(*) FROM post) AS P, " +
		"(SELECT COUNT(*) FROM users) AS U;").
		Scan(&status.Forum, &status.Thread, &status.Post, &status.User)
	return
}

func (f forumStorage) Clear() (err error) {
	_, err = f.db.Exec("TRUNCATE forum, thread, post, users, forum_users, vote CASCADE")
	return
}
