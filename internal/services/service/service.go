package service

import (
	"github.com/keithzetterstrom/db_forum/internal/models"
	"github.com/keithzetterstrom/db_forum/internal/storages"
)

type Service interface {
	ClearData() (err error)
	ReturnStatus() (status models.Status, err error)
}

type service struct {
	forumStorage storages.ForumStorage
}

func NewService(forumStorage storages.ForumStorage) Service {
	return &service{
		forumStorage: forumStorage,
	}
}

func (s service) ClearData() (err error) {
	return s.forumStorage.Clear()
}

func (s service) ReturnStatus() (status models.Status, err error) {
	return s.forumStorage.Status()
}