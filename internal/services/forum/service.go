package forum

import (
	"github.com/keithzetterstrom/db_forum/internal/models"
)

type Service interface {
	CreateForum(forumInput models.Forum) (forum models.Forum, err error)
	GetForum(slug string) (forum models.Forum, err error)
	GetForumThreads(slug string, params models.ForumQueryParams) (threads []models.Thread, err error)
	GetForumUsers(slug string, params models.ForumQueryParams) (users []models.User, err error)

	CreateThread(threadInput models.Thread) (thread models.Thread, err error)
	GetThread(slagOrID models.ThreadSlagOrID) (thread models.Thread, err error)
	UpdateThread(threadInput models.ThreadUpdate) (thread models.Thread, err error)
	GetThreadPosts(params models.ThreadQueryParams) (threads []models.Post, err error)

	ThreadVote(voteInput models.Vote) (thread models.Thread, err error)

	CreatePost(slagOrID models.ThreadSlagOrID, postInput []models.Post) (posts []models.Post, err error)
	GetPost(id int64, related []string) (post models.PostFull, err error)
	UpdatePost(postInput models.PostUpdate) (post models.Post, err error)
}

type service struct {
}

func NewService() *service {
	return &service{
	}
}

func (s *service) CreateForum(forumInput models.Forum) (forum models.Forum, err error) {
	return models.Forum{}, nil
}

func (s *service) GetForum(slug string) (forum models.Forum, err error) {
	return models.Forum{}, nil
}

func (s *service) GetForumThreads(slug string, params models.ForumQueryParams) (threads []models.Thread, err error) {
	return nil, nil
}

func (s *service) GetForumUsers(slug string, params models.ForumQueryParams) (threads []models.User, err error) {
	return nil, nil
}


func (s *service) CreateThread(threadInput models.Thread) (thread models.Thread, err error) {
	return models.Thread{}, nil
}

func (s *service) GetThread(slagOrID models.ThreadSlagOrID) (thread models.Thread, err error) {
	return models.Thread{}, nil
}

func (s *service) UpdateThread(threadInput models.ThreadUpdate) (thread models.Thread, err error) {
	return models.Thread{}, nil
}

func (s *service) GetThreadPosts(slagOrID models.ThreadSlagOrID, params models.ThreadQueryParams) (threads []models.Post, err error) {
	return nil, nil
}


func (s *service) ThreadVote(voteInput models.Vote) (thread models.Thread, err error) {
	return models.Thread{}, nil
}


func (s *service) CreatePost(slagOrID models.ThreadSlagOrID, postInput []models.Post) (posts []models.Post, err error) {
	return nil, nil
}

func (s *service) GetPost(id int64, related []string) (post models.Post, err error) {
	return models.Post{}, nil
}

func (s *service) UpdatePost(postInput models.PostUpdate) (post models.Post, err error) {
	return models.Post{}, nil
}
