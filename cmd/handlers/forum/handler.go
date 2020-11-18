package forumHandler

import "github.com/labstack/echo"

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

}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) ForumCreate(c echo.Context) error {
	return nil
}

func (h *handler) ForumGet(c echo.Context) error {
	return nil
}

func (h *handler) ForumThreadsGet(c echo.Context) error {
	return nil
}

func (h *handler) ForumUsersGet(c echo.Context) error {
	return nil
}


func (h *handler) ThreadCreate(c echo.Context) error {
	return nil
}

func (h *handler) ThreadGet(c echo.Context) error {
	return nil
}

func (h *handler) ThreadUpdate(c echo.Context) error {
	return nil
}

func (h *handler) ThreadPostsGet(c echo.Context) error {
	return nil
}


func (h *handler) ThreadVote(c echo.Context) error {
	return nil
}


func (h *handler) PostCreate(c echo.Context) error {
	return nil
}

func (h *handler) PostGet(c echo.Context) error {
	return nil
}

func (h *handler) PostUpdate(c echo.Context) error {
	return nil
}
