package profileHandler

import (
	"github.com/keithzetterstrom/db_forum/internal/models"
	"github.com/keithzetterstrom/db_forum/internal/services/profile"
	"github.com/labstack/echo"
	"net/http"
)

type Handler interface {
	ProfileCreate(c echo.Context) error
	ProfileGet(c echo.Context) error
	ProfileUpdate(c echo.Context) error
}

type handler struct {
	profileService profile.Service
}

func NewHandler(profileService profile.Service) *handler {
	return &handler{
		profileService: profileService,
	}
}

func (h *handler) ProfileCreate(c echo.Context) error {
	userInput := new(models.User)
	if err := c.Bind(userInput); err != nil {
		return err
	}

	userInput.Nickname = c.Param("nickname")

	user, err := h.profileService.CreateUser(*userInput)
	if err != nil {
		return c.JSON(http.StatusConflict, err)
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *handler) ProfileGet(c echo.Context) error {

	nickname := c.Param("nickname")

	user, err := h.profileService.GetUser(nickname)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *handler) ProfileUpdate(c echo.Context) error {
	userInput := new(models.User)
	if err := c.Bind(userInput); err != nil {
		return err
	}

	userInput.Nickname = c.Param("nickname")

	user, err := h.profileService.UpdateUser(*userInput)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, user)
}
