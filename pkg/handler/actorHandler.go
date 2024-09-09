package handler

import (
	"film_library/pkg/repository"
	"net/http"
)

type ActorHandler struct {
	repo *repository.Repository
}

func NewActorHandler(repo *repository.Repository) *ActorHandler {
	return &ActorHandler{repo: repo}
}

func (a ActorHandler) Add(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a ActorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a ActorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
