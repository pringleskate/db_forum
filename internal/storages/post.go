package storages

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
	"strconv"
	"strings"
)

type PostStorage interface {
	CreatePosts(thread models.Thread, forum string, created string, posts []models.Post) (post []models.Post, err error)

	InsertSomePosts(inputPosts []models.Post) (posts []models.Post, err error)
	InsertPost(postInput models.Post) (ID int, err error)
	UpdatePost(postInput models.PostUpdate) (err error)
	GetFullPost(postID int) (post models.Post, err error)
	GetAllPostsByThread(params models.ThreadQueryParams) (posts []models.Post, err error)
	GetPostThread(postID int) (threadID int, err error)
}

type postStorage struct {
	db *pgx.ConnPool
}

func NewPostStorage(db *pgx.ConnPool) PostStorage{
	return &postStorage{
		db: db,
	}
}

func (p postStorage) CreatePosts(thread models.Thread, forum string, created string, posts []models.Post) (post []models.Post, err error) {
	sqlStr := "INSERT INTO post(id, parent, thread, forum, author, created, message, path) VALUES "
	vals := []interface{}{}
	for _, post := range posts {
		var author string
		err = p.db.QueryRow(`SELECT nick_name FROM users WHERE lower(nick_name) = lower($1)`,
			post.Author,
		).Scan(&author)
		if err != nil {
			return nil, errors.New("404")
			//return nil, err
		}
		sqlQuery := `
		INSERT INTO public.forum_users (forum, user_nick)
		VALUES ($1,$2)`
		_, _ = p.db.Exec(sqlQuery, forum, author)

		if post.Parent == 0 {
			sqlStr += "(nextval('public.post_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"ARRAY[currval(pg_get_serial_sequence('public.post', 'id'))::INTEGER]),"
			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, created, post.Message)
		} else {
			var parentThreadId int32
			err = p.db.QueryRow("SELECT post.thread FROM post WHERE id = $1",
				post.Parent,
			).Scan(&parentThreadId)
			if err != nil {
				return nil, err
			}
			if parentThreadId != int32(thread.ID) {
				return nil, errors.New("Parent post was created in another thread")
			}

			sqlStr += " (nextval('public.post_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"(SELECT post.path FROM post WHERE post.id = ? AND post.thread = ?) || " +
				"currval(pg_get_serial_sequence('public.post', 'id'))::INTEGER),"

			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, created, post.Message, post.Parent, thread.ID)
		}

	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr += " RETURNING  id, parent, thread, forum, author, created, message, is_edited "

	sqlStr = ReplaceSQL(sqlStr, "?")
	if len(posts) > 0 {
		rows, err := p.db.Query(sqlStr, vals...)
		if err != nil {
			return nil, err
		}
		scanPost := models.Post{}
		for rows.Next() {
			err := rows.Scan(
				&scanPost.ID,
				&scanPost.Parent,
				&scanPost.ThreadID,
				&scanPost.Forum,
				&scanPost.Author,
				&scanPost.Created,
				&scanPost.Message,
				&scanPost.IsEdited,
			)
			if err != nil {
				rows.Close()
				return nil, err
			}
			post = append(post, scanPost)
		}
		rows.Close()
	}
	return post, nil
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func (p postStorage) InsertSomePosts(inputPosts []models.Post) (posts []models.Post, err error) {
	var ID int
	for i, post := range inputPosts {
		ID, err = p.InsertPost(post)
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
	_, err = p.db.Exec("UPDATE post SET message = $1, is_edited = true WHERE ID = $2", postInput.Message, postInput.ID)
	return
}

func (p postStorage) GetFullPost(postID int) (post models.Post, err error) {
	err = p.db.QueryRow("SELECT ID, author, created, forum, is_edited, message, parent, thread FROM post WHERE ID = $1", postID).
		Scan(&post.ID, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Parent, &post.ThreadID)
	return
}

func (p postStorage) GetAllPostsByThread(params models.ThreadQueryParams) (posts []models.Post, err error) {
	posts = make([]models.Post, 0)
	var rows *pgx.Rows

	switch params.Sort {
	case "":
		fallthrough
	case "flat":
		if params.Desc {
			if params.Since == 0 {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 " +
					"ORDER BY p.created DESC, p.id DESC LIMIT $2",
					params.ThreadID, params.Limit)
			} else {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 and p.id < $2 " +
					"ORDER BY p.created DESC, p.id DESC LIMIT $3",
					params.ThreadID, params.Since, params.Limit)
			}
		} else {
			if params.Since == 0 {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 " +
					"ORDER BY p.created, p.id LIMIT $2",
					params.ThreadID, params.Limit)
			} else {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 and p.id > $2 " +
					"ORDER BY p.created, p.id LIMIT $3",
					params.ThreadID, params.Since, params.Limit)
			}
		}
	case "tree":
		/*
		conditionSign := ">"
			if desc == "desc" {
				conditionSign = "<"
		}

		orderString := fmt.Sprintf(" ORDER BY p.path[1] %s, p.path %s ", desc, desc)
		sqlQuery = "SELECT p.id, p.parent, p.thread, p.forum, p.author, p.created, p.message, p.is_edited, p.path " +
			"FROM public.post as p " +
			"WHERE p.thread = $1 "
		if since != "" {
			sqlQuery += fmt.Sprintf(" AND p.path %s (SELECT p.path FROM public.post as p WHERE p.id = %s) ", conditionSign, since)
		}
		sqlQuery += orderString
		sqlQuery += fmt.Sprintf("LIMIT %s", limit)
		*/
	//	if params.Desc {
			conditionSign := ">"

			var desc string
			if params.Desc {
				conditionSign = "<"
				desc = "desc"
			} else {
				desc = ""
			}

			orderString := fmt.Sprintf(" ORDER BY p.path[1] %s, p.path %s ", desc, desc)
			sqlQuery := "SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
				"FROM post as p " +
				"WHERE p.thread = $1 "
			if params.Since != 0 {
				sqlQuery += fmt.Sprintf(" AND p.path %s (SELECT p.path FROM post as p WHERE p.id = %s) ", conditionSign, strconv.Itoa(params.Since))
			}
			sqlQuery += orderString
			sqlQuery += fmt.Sprintf("LIMIT %s", strconv.Itoa(params.Limit))

			/*orderString := fmt.Sprintf(" ORDER BY p.path[1] %s, p.path %s ", "desc", "desc")
			sqlQuery := "SELECT p.id, p.parent, p.thread, p.forum, p.author, p.created, p.message, p.is_edited, p.path " +
				"FROM public.post as p " +
				"WHERE p.thread = $1 "
			if params.Since != 0 {
				sqlQuery += fmt.Sprintf(" AND p.path %s (SELECT p.path FROM public.post as p WHERE p.id = %s) ", conditionSign, params.Since)
			}
			sqlQuery += orderString
			sqlQuery += fmt.Sprintf("LIMIT %s", params.Limit)*/

			rows, err = p.db.Query(sqlQuery, params.ThreadID)
			/*if params.Since == 0 {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1" +
					"ORDER BY p.path DESC LIMIT $2",
					params.ThreadID, params.Limit)
			} else {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 and (p.path < (SELECT p2.path from post p2 where p2.id = $2)) " +
					"ORDER BY p.path DESC LIMIT $3",
					params.ThreadID, params.Since, params.Limit)
			}*/
			/*} else {
				if params.Since == 0 {
					rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
						"FROM post p WHERE p.thread = $1" +
						"ORDER BY p.path LIMIT $2",
						params.ThreadID, params.Limit)
				} else {
					rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
						"FROM post p WHERE p.thread = $1 and (p.path > (SELECT p2.path from post p2 where p2.id = $2)) " +
						"ORDER BY p.path LIMIT $3",
						params.ThreadID, params.Since, params.Limit)
				}
			}*/
	case "parent_tree":
		if params.Desc {
			if params.Since != 0 {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 and p.path[1] IN " +
					"(SELECT p2.path[1] FROM post p2 WHERE p2.thread = $2 AND p2.parent = 0 and p2.path[1] < (SELECT p3.path[1] from post p3 where p3.id = $3) " +
					"ORDER BY p2.path DESC LIMIT $4) ORDER BY p.path[1] DESC, p.path[2:]",
					params.ThreadID, params.ThreadID, params.Since, params.Limit)
			} else {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 and p.path[1] IN " +
					"(SELECT p2.path[1] FROM post p2 WHERE p2.parent = 0 and p2.thread = $2 ORDER BY p2.path DESC LIMIT $3) " +
					"ORDER BY p.path[1] DESC, p.path[2:]",
					params.ThreadID, params.ThreadID, params.Limit)
			}
		} else {
			if params.Since != 0 {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 and p.path[1] IN " +
					"(SELECT p2.path[1] FROM post p2 WHERE p2.thread = $2 AND p2.parent = 0 and p2.path[1] > (SELECT p3.path[1] from post p3 where p3.id = $3) " +
					"ORDER BY p2.path LIMIT $4) ORDER BY p.path",
					params.ThreadID, params.ThreadID, params.Since, params.Limit)
			} else {
				rows, err = p.db.Query("SELECT p.id, p.author, p.created, p.is_edited, p.message, p.parent, p.thread, p.forum " +
					"FROM post p WHERE p.thread = $1 and p.path[1] IN " +
					"(SELECT p2.path[1] FROM post p2 WHERE p2.parent = 0 and p2.thread = $2 ORDER BY p2.path LIMIT $3) " +
					"ORDER BY path",
					params.ThreadID, params.ThreadID, params.Limit)
			}
		}
	}

	if err != nil {
		fmt.Println("rows err", err)
		return
	}
	if rows == nil {
		fmt.Println("rows nil", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited, &post.Message, &post.Parent, &post.ThreadID, &post.Forum)
		if err != nil {
			return
		}

		posts = append(posts, post)
	}

	return
}

func (p postStorage) GetPostThread(postID int) (threadID int, err error) {
	err = p.db.QueryRow("SELECT thread FROM post WHERE ID = $1", postID).Scan(&threadID)
	return
}