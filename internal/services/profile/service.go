package profile

import (
	"github.com/jackc/pgx"
	"github.com/keithzetterstrom/db_forum/internal/models"
	"github.com/keithzetterstrom/db_forum/internal/storages"
	"log"
	"strings"
)

type Service interface {
	CreateUser(userInput models.User) (user []models.User, err error)
	GetUser(nickname string) (user models.User, err error)
	UpdateUser(userInput models.User) (user models.User, err error)
}

type service struct {
	userStorage storages.UserStorage
}

func NewService(userStorage storages.UserStorage) Service {
	return &service{
		userStorage: userStorage,
	}
}

func (s service) CreateUser(userInput models.User) (user []models.User, err error) {
	user = make([]models.User, 0)

	usr, err := s.userStorage.GetFullUserByNickname(userInput.Nickname)
	if err == nil {
		user = append(user, usr)
		if strings.ToLower(usr.Email) == strings.ToLower(userInput.Email) {
			return user, models.ServError{Code: models.ConflictData}
		}
	}

	if err != pgx.ErrNoRows && err != nil {
		log.Print(err)
		return user, models.ServError{Code: models.InternalServerError}
	}

	usr, err = s.userStorage.GetFullUserByEmail(userInput.Email)
	if err == nil {
		user = append(user, usr)
	}
	if err != pgx.ErrNoRows && err != nil{
		log.Print(err)
		return user, models.ServError{Code: models.InternalServerError}
	}

	if len(user) != 0 {
		return user, models.ServError{Code: models.ConflictData}
	}

	err = s.userStorage.InsertUser(userInput)
	if err != nil {
		log.Print(err)
		return user, models.ServError{Code: models.InternalServerError}
	}

	user = append(user, userInput)
	return user, nil
}

func (s service) GetUser(nickname string) (user models.User, err error) {
	user, err = s.userStorage.GetFullUserByNickname(nickname)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return user, models.ServError{Code: models.InternalServerError}
	}

	return user, nil
}

func (s service) UpdateUser(userInput models.User) (user models.User, err error) {
	user, err = s.userStorage.GetFullUserByNickname(userInput.Nickname)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, models.ServError{Code: models.NotFound}
		}
		log.Print(err)
		return user, models.ServError{Code: models.InternalServerError}
	}

	_, err = s.userStorage.GetFullUserByEmail(userInput.Email)
	if err == nil {
		return user, models.ServError{Code: models.ConflictData}
	}
	if err != pgx.ErrNoRows {
		log.Print(err)
		return user, models.ServError{Code: models.InternalServerError}
	}

	if userInput.Email != "" {
		user.Email = userInput.Email
	}
	if userInput.About != "" {
		user.About = userInput.About
	}
	if userInput.FullName != "" {
		user.FullName = userInput.FullName
	}

	err = s.userStorage.UpdateUser(user)
	if err != nil {
		log.Print(err)
		return user, models.ServError{Code: models.InternalServerError}
	}

	return user, nil
}