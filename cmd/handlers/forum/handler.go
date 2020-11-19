package forumHandler

import (
	"github.com/keithzetterstrom/db_forum/internal/models"
	"github.com/keithzetterstrom/db_forum/internal/services/forum"
	"github.com/labstack/echo"
	"net/url"
	"strconv"
)

type Handler interface {
	ForumCreate(c echo.Context) error
	ForumGet(c echo.Context) error
	ForumThreadsGet(c echo.Context) error
	ForumUsersGet(c echo.Context) error

	ThreadCreate(c echo.Context) error
	ThreadGet(c echo.Context) error
	ThreadUpdate(c echo.Context) error
	ThreadPostsGet(c echo.Context) error

	ThreadVote(c echo.Context) error

	PostCreate(c echo.Context) error
	PostGet(c echo.Context) error
	PostUpdate(c echo.Context) error
}

type handler struct {
	forumService forum.Service
}

func NewHandler(forumService forum.Service) *handler {
	return &handler{
		forumService: forumService,
	}
}

func (h *handler) ForumCreate(c echo.Context) error {
	forumInput := new(models.Forum)
	if err := c.Bind(forumInput); err != nil {
		return err
	}

	_, err := h.forumService.CreateForum(*forumInput)
	if err != nil {

	}

	return nil
}

func (h *handler) ForumGet(c echo.Context) error {
	slag := c.Param("slug")

	_, err := h.forumService.GetForum(slag)
	if err != nil {

	}

	return nil
}

func (h *handler) ForumThreadsGet(c echo.Context) error {
	slag := c.Param("slug")
	params, err := getForumQueryParams(c.QueryParams())
	if err != nil {
		return err
	}

	_, err = h.forumService.GetForumThreads(slag, params)
	if err != nil {

	}

	return nil
}

func (h *handler) ForumUsersGet(c echo.Context) error {
	slag := c.Param("slug")
	params, err := getForumQueryParams(c.QueryParams())
	if err != nil {
		return err
	}

	_, err = h.forumService.GetForumUsers(slag, params)
	if err != nil {

	}

	return nil
}


func (h *handler) ThreadCreate(c echo.Context) error {
	threadInput := new(models.Thread)
	if err := c.Bind(threadInput); err != nil {
		return err
	}

	threadInput.ForumSlug = c.Param("slug")

	_, err := h.forumService.CreateThread(*threadInput)
	if err != nil {

	}

	return nil
}

func (h *handler) ThreadGet(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	_, err := h.forumService.GetThread(slugOrID)
	if err != nil {

	}

	return nil
}

func (h *handler) ThreadUpdate(c echo.Context) error {
	threadInput := new(models.ThreadUpdate)
	if err := c.Bind(threadInput); err != nil {
		return err
	}

	threadInput.SlagOrID = c.Param("slug_or_id")

	_, err := h.forumService.UpdateThread(*threadInput)
	if err != nil {

	}

	return nil
}

func (h *handler) ThreadPostsGet(c echo.Context) error {
	slagOrID := c.Param("slug_or_id")
	params, err := getThreadQueryParams(c.QueryParams())
	if err != nil {
		return err
	}

	_, err = h.forumService.GetThreadPosts(slagOrID, params)
	if err != nil {

	}

	return nil
}


func (h *handler) ThreadVote(c echo.Context) error {
	voteInput := new(models.Vote)
	if err := c.Bind(voteInput); err != nil {
		return err
	}

	voteInput.SlagOrID = c.Param("slug_or_id")

	_, err := h.forumService.ThreadVote(*voteInput)
	if err != nil {

	}

	return nil
}


func (h *handler) PostCreate(c echo.Context) error {
	//TODO - приходит список создаваемых постов
	postInput := make([]models.Post, 0)
	if err := c.Bind(postInput); err != nil {
		return err
	}

	slagOrID := c.Param("slug_or_id")

	_, err := h.forumService.CreatePost(slagOrID, postInput)
	if err != nil {

	}

	return nil
}

func (h *handler) PostGet(c echo.Context) error {
	//TODO - есть параметр очернеди related (array [string])
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	_, err = h.forumService.GetPost(id)
	if err != nil {

	}

	return nil
}

func (h *handler) PostUpdate(c echo.Context) (err error) {
	postInput := new(models.PostUpdate)
	if err := c.Bind(postInput); err != nil {
		return err
	}

	postInput.ID , err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	_, err = h.forumService.UpdatePost(*postInput)
	if err != nil {

	}

	return nil
}

func getForumQueryParams(params url.Values) (values models.ForumQueryParams, err error) {
	values.Limit, err = strconv.Atoi(params.Get("limit"))
	values.Since = params.Get("since")
	values.Desc, err = strconv.ParseBool(params.Get("desc"))
	if err != nil {
		return models.ForumQueryParams{}, err
	}
	return values, nil
}

func getThreadQueryParams(params url.Values) (values models.ThreadQueryParams, err error) {
	values.Limit, err = strconv.Atoi(params.Get("limit"))
	values.Since = params.Get("since")
	values.Sort = params.Get("sort")
	values.Desc, err = strconv.ParseBool(params.Get("desc"))
	if err != nil {
		return models.ThreadQueryParams{}, err
	}
	return values, nil
}