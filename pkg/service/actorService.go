package service

import (
	"film_library/model"
	"film_library/pkg/repository"
)

type ActorService struct {
	repo *repository.Repository
}

func NewActorService(repo *repository.Repository) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) Add(actorId int, actor *model.Actor) (*model.Actor, error) {
	return s.repo.Add(actorId, actor)
}

func (s *ActorService) Update(actorId int, actor *model.Actor) (*model.Actor, error) {
	return s.repo.Update(actorId, actor)
}

func (s *ActorService) Delete(actorId int) (*model.Actor, error) {
	return s.repo.Delete(actorId)
}
