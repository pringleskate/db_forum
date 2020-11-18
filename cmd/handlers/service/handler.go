package serviceHandler

import "github.com/labstack/echo"

type Handler interface {
	ServiceClear(c echo.Context) error
	ServiceStatus(c echo.Context) error
}

type handler struct {

}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) ServiceClear(c echo.Context) error {
	return nil
}

func (h *handler) ServiceStatus(c echo.Context) error {
	return nil
}