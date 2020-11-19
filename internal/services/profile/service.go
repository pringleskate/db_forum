package profile

import "github.com/keithzetterstrom/db_forum/internal/models"

type Service interface {
	CreateUser(userInput models.User) (user models.User, err error)
	GetUser(nickname string) (user models.User, err error)
	UpdateUser(userInput models.UsersUpdate) (user models.User, err error)
}

type service struct {}

func NewService() Service {
	return &service{}
}

func (s service) CreateUser(userInput models.User) (user models.User, err error) {
	return user, nil
}

func (s service) GetUser(nickname string) (user models.User, err error) {
	return user, nil
}

func (s service) UpdateUser(userInput models.UsersUpdate) (user models.User, err error) {
	return user, nil
}