package forum

import (
	"github.com/keithzetterstrom/db_forum/internal/models"
)

type Service interface {
	CreateForum(forumInput models.Forum) (forum models.Forum, err error)
	GetForum(slug string) (forum models.Forum, err error)
	GetForumThreads(slug string, params models.ForumQueryParams) (threads []models.Thread, err error)
	GetForumUsers(slug string, params models.ForumQueryParams) (threads []models.User, err error)

	CreateThread(threadInput models.Thread) (thread models.Thread, err error)
	GetThread(slagOrID string) (thread models.Thread, err error)
	UpdateThread(threadInput models.ThreadUpdate) (thread models.Thread, err error)
	GetThreadPosts(slagOrID string, params models.ThreadQueryParams) (threads []models.Post, err error)

	ThreadVote(voteInput models.Vote) (thread models.Thread, err error)

	CreatePost(slagOrID string, postInput []models.Post) (posts []models.Post, err error)
	GetPost(id int64) (post models.Post, err error)
	UpdatePost(postInput models.PostUpdate) (post models.Post, err error)
}

type service struct {

}

func NewService() *service {
	return &service{
	}
}
