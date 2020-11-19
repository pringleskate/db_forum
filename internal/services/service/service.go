package service

import "github.com/keithzetterstrom/db_forum/internal/models"

type Service interface {
	ClearData() (err error)
	ReturnStatus() (status models.Status, err error)
}

type service struct {}

func NewService() Service {
	return &service{}
}

func (s service) ClearData() (err error) {
	return nil
}

func (s service) ReturnStatus() (status models.Status, err error) {
	return status, nil
}