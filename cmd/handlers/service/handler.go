package serviceHandler

import (
	"github.com/keithzetterstrom/db_forum/internal/services/service"
	"github.com/labstack/echo"
	"net/http"
)

type Handler interface {
	ServiceClear(c echo.Context) error
	ServiceStatus(c echo.Context) error
}

type handler struct {
	serviceService service.Service
}

func NewHandler(serviceService service.Service) *handler {
	return &handler{
		serviceService: serviceService,
	}
}

func (h *handler) ServiceClear(c echo.Context) error {
	err := h.serviceService.ClearData()
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) ServiceStatus(c echo.Context) error {
	status, err := h.serviceService.ReturnStatus()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, status)
}