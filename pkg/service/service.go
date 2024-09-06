package service

import (
	"film_library/model"
	"film_library/pkg/repository"
)

type Service struct {
	actorService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{actorService: NewActorService(repo)}
}

type actorService interface {
	Add(actorId int, actor *model.Actor) (*model.Actor, error)
	Update(actorId int, actor *model.Actor) (*model.Actor, error)
	Delete(actorId int) (*model.Actor, error)
}
