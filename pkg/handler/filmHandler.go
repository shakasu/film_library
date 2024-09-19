package handler

import (
	"film_library/pkg/repository"
	"net/http"
)

type FilmHandler struct {
	repo *repository.Repository
}

func NewFilmHandler(repo *repository.Repository) *FilmHandler {
	return &FilmHandler{repo: repo}
}

func (f FilmHandler) Add(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (f FilmHandler) Update(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (f FilmHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (f FilmHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (f FilmHandler) GetByFilter(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
