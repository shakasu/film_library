package handler

import (
	"encoding/json"
	"errors"
	"film_library/model"
	"film_library/pkg/repository"
	"log"
	"net/http"
)

type FilmHandler struct {
	repo *repository.Repository
}

func NewFilmHandler(repo *repository.Repository) *FilmHandler {
	return &FilmHandler{repo: repo}
}

func (handler FilmHandler) Add(w http.ResponseWriter, r *http.Request) {
	if !authWriter(w, r, handler.repo) {
		return
	}
	var filmReq model.FilmDto

	err := decodeJSONBody(w, r, &filmReq)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	filmResp, err := handler.repo.FilmRepo.Add(&filmReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(filmResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		log.Print(err.Error())
	}
}

func (handler FilmHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !authWriter(w, r, handler.repo) {
		return
	}
	var filmReq model.FilmDto

	id, err := getIdFromPath(w, r)
	if err != nil {
		return
	}

	err = decodeJSONBody(w, r, &filmReq)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	filmResp, err := handler.repo.FilmRepo.Update(&filmReq, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(filmResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		log.Print(err.Error())
	}
}

func (handler FilmHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !authWriter(w, r, handler.repo) {
		return
	}
	id, err := getIdFromPath(w, r)
	if err != nil {
		return
	}
	err = handler.repo.FilmRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler FilmHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if !authReader(w, r, handler.repo) {
		return
	}
	actors, err := handler.repo.FilmRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(actors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		log.Print(err.Error())
	}
}
