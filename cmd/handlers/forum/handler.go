package forumHandler

import (
	"github.com/keithzetterstrom/db_forum/internal/models"
	"github.com/keithzetterstrom/db_forum/internal/services/forum"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	forumRequest, err := h.forumService.CreateForum(*forumInput)
	if err != nil {

	}

	return c.JSON(http.StatusCreated, forumRequest)
}

func (h *handler) ForumGet(c echo.Context) error {
	slag := c.Param("slug")

	forumRequest, err := h.forumService.GetForum(slag)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, forumRequest)
}

func (h *handler) ForumThreadsGet(c echo.Context) error {
	slag := c.Param("slug")
	params, err := getForumQueryParams(c.QueryParams())
	if err != nil {
		return err
	}

	threads, err := h.forumService.GetForumThreads(slag, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, threads)
}

func (h *handler) ForumUsersGet(c echo.Context) error {
	slag := c.Param("slug")
	params, err := getForumQueryParams(c.QueryParams())
	if err != nil {
		return err
	}

	users, err := h.forumService.GetForumUsers(slag, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}


func (h *handler) ThreadCreate(c echo.Context) error {
	threadInput := new(models.Thread)
	if err := c.Bind(threadInput); err != nil {
		return err
	}

	threadInput.ForumSlug = c.Param("slug")

	thread, err := h.forumService.CreateThread(*threadInput)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, thread)
}

func (h *handler) ThreadGet(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	thread, err := h.forumService.GetThread(slugOrID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, thread)
}

func (h *handler) ThreadUpdate(c echo.Context) error {
	threadInput := new(models.ThreadUpdate)
	if err := c.Bind(threadInput); err != nil {
		return err
	}

	threadInput.ThreadSlagOrID = isItSlugOrID(c.Param("slug_or_id"))

	thread, err := h.forumService.UpdateThread(*threadInput)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, thread)
}

func (h *handler) ThreadPostsGet(c echo.Context) error {
	params, err := getThreadQueryParams(c.QueryParams())
	if err != nil {
		return err
	}

	params.ThreadSlagOrID = isItSlugOrID(c.Param("slug_or_id"))

	posts, err := h.forumService.GetThreadPosts(params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, posts)
}


func (h *handler) ThreadVote(c echo.Context) error {
	voteInput := new(models.Vote)
	if err := c.Bind(voteInput); err != nil {
		return err
	}

	voteInput.SlagOrID = c.Param("slug_or_id")

	thread, err := h.forumService.ThreadVote(*voteInput)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, thread)
}


func (h *handler) PostCreate(c echo.Context) error {
	postInput := make([]models.Post, 0)
	if err := c.Bind(postInput); err != nil {
		return err
	}

	slagOrID := isItSlugOrID(c.Param("slug_or_id"))

	post, err := h.forumService.CreatePost(slagOrID, postInput)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)
}

func (h *handler) PostGet(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	related := relatedParse(c.QueryParam("related"))

	post, err := h.forumService.GetPost(id, related)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)
}

func (h *handler) PostUpdate(c echo.Context) (err error) {
	postInput := new(models.PostUpdate)
	if err := c.Bind(postInput); err != nil {
		return err
	}

	postInput.ID , err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	post, err := h.forumService.UpdatePost(*postInput)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)
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

func relatedParse(related string) []string {
	related = strings.ReplaceAll(related, "[", "")
	related = strings.ReplaceAll(related, "]", "")
	return strings.Split(related, ",")
}

func isItSlugOrID(slagOrID string) (output models.ThreadSlagOrID) {
	id, err := strconv.Atoi(slagOrID)
	if err != nil {
		output.ThreadSlug = slagOrID
		return output
	}
	output.ThreadID = id
	return output
}