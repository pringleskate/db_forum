package forum

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
	"github.com/keithzetterstrom/db_forum/internal/storages"
	"log"
	"strings"
	"time"
)

type Service interface {
	CreateForum(forumInput models.Forum) (forum models.Forum, err error)
	GetForum(slug string) (forum models.Forum, err error)
	GetForumThreads(slug string, params models.ForumQueryParams) (threads []models.Thread, err error)
	GetForumUsers(slug string, params models.ForumQueryParams) (users []models.User, err error)

	CreateThread(threadInput models.Thread) (thread models.Thread, err error)
	GetThread(slagOrID models.ThreadSlagOrID) (thread models.Thread, err error)
	UpdateThread(threadInput models.ThreadUpdate) (thread models.Thread, err error)
	GetThreadPosts(params models.ThreadQueryParams) (posts []models.Post, err error)

	ThreadVote(voteInput models.Vote) (thread models.Thread, err error)

	CreatePosts(slagOrID models.ThreadSlagOrID, postInput []models.Post) (posts []models.Post, err error)
	GetPost(id int64, related []string) (post models.PostFull, err error)
	UpdatePost(postInput models.PostUpdate) (post models.Post, err error)
}

type service struct {
	forumStorage storages.ForumStorage
	threadStorage storages.ThreadStorage
	postStorage storages.PostStorage
	userStorage storages.UserStorage
}

func NewService(forumStorage storages.ForumStorage, threadStorage storages.ThreadStorage, postStorage storages.PostStorage, userStorage storages.UserStorage) Service {
//func NewService(forumStorage storages.ForumStorage, threadStorage storages.ThreadStorage, postStorage storages.PostStorage, userStorage storages.UserStorage) *service {
	return &service{
		forumStorage: forumStorage,
		threadStorage: threadStorage,
		postStorage: postStorage,
		userStorage: userStorage,
	}
}

func (s *service) CreateForum(forumInput models.Forum) (forum models.Forum, err error) {
	existing, err := s.forumStorage.GetFullForum(forumInput.Slug)
	if err == nil {
		return existing, models.ServError{ Code: models.ConflictData}
	}
	if err != pgx.ErrNoRows {
		log.Print(err)
		return forum, models.ServError{ Code: models.InternalServerError }
	}

	nickname, err := s.userStorage.GetUserNickname(forumInput.User)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forum, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return forum, models.ServError{ Code: models.InternalServerError }
	}

	forumInput.User = nickname

	err = s.forumStorage.InsertForum(forumInput)
	if err != nil {
		log.Print(err)
		return forum, models.ServError{ Code: models.InternalServerError }
	}

	return forumInput, nil
}

func (s *service) GetForum(slug string) (forum models.Forum, err error) {
	forum, err = s.forumStorage.GetFullForum(slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forum, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return forum, models.ServError{ Code: models.InternalServerError }
	}

	return forum, nil
}

func (s *service) GetForumThreads(slug string, params models.ForumQueryParams) (threads []models.Thread, err error) {
	threads = make([]models.Thread, 0)

	forum, err := s.forumStorage.GetFullForum(slug)
	if err != nil {
		log.Print(err)
		if err == pgx.ErrNoRows {
			return threads, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return threads, models.ServError{ Code: models.InternalServerError }
	}

	params.Slug = forum.Slug
	if params.Limit == 0 {
		params.Limit = 10000
	}

	threads, err = s.threadStorage.GetAllThreadsByForum(params)
	if err != nil {
		log.Print(err)
		return threads, models.ServError{ Code: models.InternalServerError }
	}

	return threads, nil
}

func (s *service) GetForumUsers(slug string, params models.ForumQueryParams) (users []models.User, err error) {
	users = make([]models.User, 0)

	forum, err := s.forumStorage.GetFullForum(slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return users, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return users, models.ServError{ Code: models.InternalServerError }
	}

	params.Slug = forum.Slug
	if params.Limit == 0 {
		params.Limit = 10000
	}

	users, err = s.userStorage.GetAllUsersByForum(params)
	if err != nil {
		log.Print(err)
		return users, models.ServError{ Code: models.InternalServerError }
	}

	return users, nil
}

func (s *service) CreateThread(threadInput models.Thread) (thread models.Thread, err error) {
	forum, err := s.forumStorage.GetFullForum(threadInput.Forum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return thread, models.ServError{Code: models.NotFound, Message: "No such forum"}
		}
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	threadInput.Forum = forum.Slug

	if threadInput.Slag != "" {
		existing, err := s.threadStorage.GetFullThreadBySlug(threadInput.Slag)
		if err == nil {
			return existing, models.ServError{ Code: models.ConflictData }
		}
		if err != pgx.ErrNoRows {
			log.Print(err)
			return thread, models.ServError{ Code: models.InternalServerError }
		}
	}

	nickname, err := s.userStorage.GetUserNickname(threadInput.Author)
	if err != nil {
		if err == pgx.ErrNoRows {
			return thread, models.ServError{Code: models.NotFound, Message: "No such user"}
		}
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}


	threadInput.Author = nickname
	ID, err := s.threadStorage.InsertThread(threadInput)
	if err != nil {
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	err = s.forumStorage.InsertForumUser(threadInput.Forum, threadInput.Author)
	if err != nil {
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	err = s.forumStorage.UpdateThreadsCount(threadInput.Forum)
	if err != nil {
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	threadInput.ID = ID

	return threadInput, nil
}

func (s *service) GetThread(slagOrID models.ThreadSlagOrID) (thread models.Thread, err error) {
	if slagOrID.ThreadSlug != "" {
		thread, err = s.threadStorage.GetFullThreadBySlug(slagOrID.ThreadSlug)
	} else {
		thread, err = s.threadStorage.GetFullThreadByID(slagOrID.ThreadID)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return thread, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	return thread, nil
}

func (s *service) UpdateThread(threadInput models.ThreadUpdate) (thread models.Thread, err error) {
	if threadInput.ThreadSlug != "" {
		thread, err = s.threadStorage.GetFullThreadBySlug(threadInput.ThreadSlug)
	} else {
		thread, err = s.threadStorage.GetFullThreadByID(threadInput.ThreadID)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return thread, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	if threadInput.Title != "" {
		thread.Title = threadInput.Title
	}
	if threadInput.Message != "" {
		thread.Message = threadInput.Message
	}
	if threadInput.Title == "" && threadInput.Message == "" {
		return thread, nil
	}

	err = s.threadStorage.UpdateThread(thread)
	if err != nil {
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	return thread, nil
}

func (s *service) GetThreadPosts(params models.ThreadQueryParams) (posts []models.Post, err error) {
	posts = make([]models.Post, 0)

	thread := models.Thread{}
	if params.ThreadSlug != "" {
		thread, err = s.threadStorage.GetFullThreadBySlug(params.ThreadSlug)
	} else {
		thread, err = s.threadStorage.GetFullThreadByID(params.ThreadID)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return posts, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

	if params.Limit == 0 {
		params.Limit = 10000
	}
	params.ThreadID = thread.ID

	posts, err = s.postStorage.GetAllPostsByThread(params)
	if err != nil {
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

	return posts, nil
}

func (s *service) ThreadVote(voteInput models.Vote) (thread models.Thread, err error) {
	if voteInput.ThreadSlug != "" {
		thread, err = s.threadStorage.GetFullThreadBySlug(voteInput.ThreadSlug)
	} else {
		thread, err = s.threadStorage.GetFullThreadByID(voteInput.ThreadID)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return thread, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}
	voteInput.ThreadID = thread.ID

	user, err := s.userStorage.GetUserNickname(voteInput.Nickname)
	if err != nil {
		if err == pgx.ErrNoRows {
			return thread, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}
	voteInput.Nickname = user

	vote, err := s.threadStorage.SelectVote(voteInput)
	if err != nil && err != pgx.ErrNoRows {
		log.Print(err)
		return thread, models.ServError{ Code: models.InternalServerError }
	}

	var newVoice int
	if err == nil {
		//update vote
		if voteInput.Voice == vote.Voice {
			return thread, nil
		}
		err = s.threadStorage.UpdateVote(voteInput)
		if err != nil {
			log.Print(err)
			return thread, models.ServError{ Code: models.InternalServerError }
		}

		if voteInput.Voice == -1 {
			newVoice = -2
		} else {
			newVoice = 2
		}

		err = s.threadStorage.UpdateVotesCount(voteInput.ThreadID, newVoice)
		if err != nil {
			log.Print(err)
			return thread, models.ServError{ Code: models.InternalServerError }
		}

	} else {
		//insert vote

		err = s.threadStorage.InsertVote(voteInput)
		if err != nil {
			log.Print(err)
			return thread, models.ServError{ Code: models.InternalServerError }
		}

		newVoice = int(voteInput.Voice)
		err = s.threadStorage.UpdateVotesCount(voteInput.ThreadID, newVoice)
		if err != nil {
			log.Print(err)
			return thread, models.ServError{ Code: models.InternalServerError }
		}
	}

	thread.Votes += newVoice
	return thread, nil
}
/*
func (s *service) CreatePosts(slagOrID models.ThreadSlagOrID, postInput []models.Post) (posts []models.Post, err error) {
	posts = make([]models.Post, 0)

	thread := models.Thread{}
	if slagOrID.ThreadSlug != "" {
		thread, err = s.threadStorage.GetFullThreadBySlug(slagOrID.ThreadSlug)
	} else {
		thread, err = s.threadStorage.GetFullThreadByID(slagOrID.ThreadID)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return posts, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

 	if len(postInput) == 0 {
		return []models.Post{}, nil
	}
	createdTime := time.Now().Format(time.RFC3339Nano)

	//createdTime, _ := time.Parse(time.RFC3339Nano, time.Now().String())
	forum := thread.Forum

	for i, post := range postInput {
		postInput[i].Forum = forum
		post.Forum = forum

		postInput[i].Created = createdTime
		post.Created = createdTime

		nickname, err := s.userStorage.GetUserNickname(post.Author)
		if err != nil {
			if err == pgx.ErrNoRows {
				return posts, models.ServError{Code: models.NotFound}
			}
			log.Print(err)
			return posts, models.ServError{ Code: models.InternalServerError }
		}

		postInput[i].Author = nickname
		post.Author = nickname

		postInput[i].ThreadID = thread.ID
		post.ThreadID = thread.ID

		postInput[i].ThreadSlug = thread.Slag
		post.ThreadSlug = thread.Slag

		if post.Parent != 0 {
			threadID, err := s.postStorage.GetPostThread(post.Parent)
			if err != nil {
				log.Print(err)
				return posts, models.ServError{ Code: models.ConflictData }
			}

			if threadID != postInput[i].ThreadID {
				return posts, models.ServError{ Code: models.ConflictData }
			}
		}

		err = s.forumStorage.InsertForumUser(post.Forum, post.Author)
		if err != nil {
			log.Print(err)
			return posts, models.ServError{ Code: models.InternalServerError }
		}

		ID, err := s.postStorage.InsertPost(postInput[i])
		if err != nil {
			log.Print(err)
			return posts, models.ServError{ Code: models.InternalServerError }
		}
		postInput[i].ID = ID
	}

	output, err := s.postStorage.InsertSomePosts(postInput)
	if err != nil {
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

err = s.forumStorage.UpdatePostsCount(forum, len(output))

	if err != nil {
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

	return output, nil
	err = s.forumStorage.UpdatePostsCount(forum, len(postInput))
	if err != nil {
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

	return postInput, nil
}*/




func (s *service) CreatePosts(slagOrID models.ThreadSlagOrID, postInput []models.Post) (posts []models.Post, err error) {
	posts = make([]models.Post, 0)

	thread := models.Thread{}
	if slagOrID.ThreadSlug != "" {
		thread, err = s.threadStorage.GetFullThreadBySlug(slagOrID.ThreadSlug)
	} else {
		thread, err = s.threadStorage.GetFullThreadByID(slagOrID.ThreadID)
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return posts, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

	if len(postInput) == 0 {
		return []models.Post{}, nil
	}

	createdTime := time.Now().Format(time.RFC3339Nano)
	//createdTime, _ := time.Parse(time.RFC3339Nano, time.Now().String())
	forum := thread.Forum

	posts, err =  s.postStorage.CreatePosts(thread, forum, createdTime, postInput)
	if err != nil {
		if err.Error() == "404" {
			return posts, models.ServError{Code: models.NotFound}
		}
		fmt.Println(err)
		return posts, models.ServError{ Code: models.ConflictData }
	}

	err = s.forumStorage.UpdatePostsCount(forum, len(posts))
	if err != nil {
		log.Print(err)
		return posts, models.ServError{ Code: models.InternalServerError }
	}

	return posts, nil
}

func (s *service) GetPost(id int64, related []string) (post models.PostFull, err error) {
	onePost, err := s.postStorage.GetFullPost(int(id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return post, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return post, models.ServError{ Code: models.InternalServerError }
	}

	post.Post = &onePost
	args := fmt.Sprint(related)
	if strings.Contains(args, "user") {
		user, err := s.userStorage.GetFullUserByNickname(post.Post.Author)
		if err != nil {
			log.Print(err)
			return post, models.ServError{ Code: models.InternalServerError }
		}
		post.Author = &user
	}
	if strings.Contains(args, "forum") {
		forum, err := s.forumStorage.GetFullForum(post.Post.Forum)
		if err != nil {
			log.Print(err)
			return post, models.ServError{ Code: models.InternalServerError }
		}
		post.Forum = &forum
	}
	if strings.Contains(args, "thread") {
		thread, err := s.threadStorage.GetFullThreadByID(post.Post.ThreadID)
		if err != nil {
			log.Print(err)
			return post, models.ServError{ Code: models.InternalServerError }
		}
		post.Thread = &thread
	}

	return post, nil
}

func (s *service) UpdatePost(postInput models.PostUpdate) (post models.Post, err error) {
	post, err = s.postStorage.GetFullPost(postInput.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return post, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return post, models.ServError{ Code: models.InternalServerError }
	}

	if postInput.Message != "" && postInput.Message != post.Message {
		err = s.postStorage.UpdatePost(postInput)
		if err != nil {
			log.Print(err)
			return post, models.ServError{ Code: models.InternalServerError }
		}

		post.Message = postInput.Message
		post.IsEdited = true
	}

	return post, nil
}
