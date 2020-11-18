package profileHandler

import "github.com/labstack/echo"

type Handler interface {
	ProfileCreate(c echo.Context) error
	ProfileGet(c echo.Context) error
	ProfileUpdate(c echo.Context) error
}

type handler struct {

}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) ProfileCreate(c echo.Context) error {
	return nil
}

func (h *handler) ProfileGet(c echo.Context) error {
	return nil
}

func (h *handler) ProfileUpdate(c echo.Context) error {
	return nil
}
