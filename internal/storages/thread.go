package storages

import (
	"database/sql"
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
)

type ThreadStorage interface {
	InsertThread(threadInput models.Thread) (ID int, err error)
	UpdateThread(threadInput models.Thread) (err error)
	GetFullThreadBySlug(slug string) (thread models.Thread, err error)
	GetFullThreadByID(threadID int) (thread models.Thread, err error)
	InsertVote(voteInput models.Vote) (err error)
	UpdateVote(voteInput models.Vote) (err error)
	SelectVote(voteInput models.Vote) (vote models.Vote, err error)
	GetAllThreadsByForum(params models.ForumQueryParams) (threads []models.Thread, err error)
	UpdateVotesCount(threadID int, voice int) (err error)
}

type threadStorage struct {
	db *pgx.ConnPool
}

func NewThreadStorage(db *pgx.ConnPool) ThreadStorage{
	return &threadStorage{
		db: db,
	}
}

func (t threadStorage) InsertThread(threadInput models.Thread) (ID int, err error) {
	err = t.db.QueryRow("INSERT INTO thread (author, created, forum, message, slug, title) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID",
						threadInput.Author, threadInput.Created, threadInput.Forum, threadInput.Message, threadInput.Slag, threadInput.Title).Scan(&ID)
	return
}

func (t threadStorage) UpdateThread(threadInput models.Thread) (err error) {
	_, err = t.db.Exec("UPDATE thread SET message = $1, title = $2 WHERE ID = $3", threadInput.Message, threadInput.Title, threadInput.ID)
	return
}

func (t threadStorage) GetFullThreadBySlug(slug string) (thread models.Thread, err error) {
	err = t.db.QueryRow("SELECT ID, author, created, forum, message, slug, title, votes FROM thread WHERE lower(slug) = lower($1)", slug).
				Scan(&thread.ID, &thread.Author, &thread.Created, &thread.Forum, &thread.Message, &thread.Slag, &thread.Title, &thread.Votes)
	return
}

func (t threadStorage) GetFullThreadByID(threadID int) (thread models.Thread, err error) {
	var slug sql.NullString
	err = t.db.QueryRow("SELECT ID, author, created, forum, message, slug, title, votes FROM thread WHERE ID = $1", threadID).
		Scan(&thread.ID, &thread.Author, &thread.Created, &thread.Forum, &thread.Message, &slug, &thread.Title, &thread.Votes)

	if slug.Valid {
		thread.Slag = slug.String
	}

	return
}

func (t threadStorage) InsertVote(voteInput models.Vote) (err error) {
	_, err = t.db.Exec("INSERT INTO vote (user_nick, voice, thread_id) VALUES ($1, $2, $3)", voteInput.Nickname, voteInput.Voice, voteInput.ThreadID)
	return
}

func (t threadStorage) UpdateVote(voteInput models.Vote) (err error) {
	_, err = t.db.Exec("UPDATE vote SET voice = $1 WHERE user_nick = $2 AND thread_id = $3", voteInput.Voice, voteInput.Nickname, voteInput.ThreadID)
	return
}

func (t threadStorage) SelectVote(voteInput models.Vote) (vote models.Vote, err error) {
	err = t.db.QueryRow("SELECT user_nick, voice, thread_id FROM vote WHERE lower(user_nick) = lower($1) AND thread_id = $2",
		voteInput.Nickname, voteInput.ThreadID).
		Scan(&vote.Nickname, &vote.Voice, &vote.ThreadID)
	return
}

func (t threadStorage) UpdateVotesCount(threadID int, voice int) (err error) {
	_, err = t.db.Exec("UPDATE thread SET votes = votes + $1 WHERE ID = $2", voice, threadID)
	return
}

func (t threadStorage) GetAllThreadsByForum(params models.ForumQueryParams) (threads []models.Thread, err error) {
	var rows *pgx.Rows
	threads = make([]models.Thread, 0)

	if params.Desc {
		if params.Since == "" {
			rows, err = t.db.Query("SELECT ID, author, created, forum, message, slug, title, votes "+
				"FROM thread WHERE lower(forum) = lower($1) ORDER BY created DESC LIMIT $2",
				params.Slug, params.Limit)
		} else {
			rows, err = t.db.Query("SELECT ID, author, created, forum, message, slug, title, votes "+
				"FROM thread WHERE lower(forum) = lower($1) AND created <= $3 ORDER BY created DESC LIMIT $2",
				params.Slug, params.Limit, params.Since)
		}
	} else {
		if params.Since == "" {
			rows, err = t.db.Query("SELECT ID, author, created, forum, message, slug, title, votes "+
				"FROM thread WHERE lower(forum) = lower($1) ORDER BY created LIMIT $2",
				params.Slug, params.Limit)
		} else {
			rows, err = t.db.Query("SELECT ID, author, created, forum, message, slug, title, votes "+
				"FROM thread WHERE lower(forum) = lower($1) AND created >= $3 ORDER BY created LIMIT $2",
				params.Slug, params.Limit, params.Since)
		}
	}
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		tmpThread := models.Thread{}
		var slug sql.NullString

		err = rows.Scan(&tmpThread.ID, &tmpThread.Author, &tmpThread.Created, &tmpThread.Forum, &tmpThread.Message, &slug, &tmpThread.Title, &tmpThread.Votes)
		if err != nil {
			return
		}

		if slug.Valid {
			tmpThread.Slag = slug.String
		}

		threads = append(threads, tmpThread)
	}

	return
}
