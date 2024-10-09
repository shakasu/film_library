package handler

import (
	"encoding/json"
	"errors"
	"film_library/model"
	"film_library/pkg/repository"
	"log"
	"net/http"
)

type ActorHandler struct {
	repo *repository.Repository
}

func NewActorHandler(repo *repository.Repository) *ActorHandler {
	return &ActorHandler{repo: repo}
}

// @Summary Создание актера
// @Tags actor
// @Param data body model.ActorDto true "The input actor struct"
// @accept  json
// @Produce  json
// @Success 200
// @Router /actor [post]
// @Security BasicAuth
func (handler ActorHandler) Add(w http.ResponseWriter, r *http.Request) {
	if !authWriter(w, r, handler.repo) {
		return
	}
	var actorReq model.ActorDto

	err := decodeJSONBody(w, r, &actorReq)
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

	actorResp, err := handler.repo.ActorRepo.Add(&actorReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(actorResp)
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

// @Summary Редактирование актера по ID
// @Tags actor
// @Param        id   path      int  true  "Actor ID"
// @Param data body model.ActorDto true "The input actor struct"
// @accept  json
// @Produce  json
// @Success 200
// @Router /actor/{id} [put]
// @Security BasicAuth
func (handler ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !authWriter(w, r, handler.repo) {
		return
	}
	var actorReq model.ActorDto

	id, err := getIdFromPath(w, r)
	if err != nil {
		return
	}

	err = decodeJSONBody(w, r, &actorReq)
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

	actorResp, err := handler.repo.ActorRepo.Update(&actorReq, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(actorResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		return
	}
}

// @Summary Удаление актера по ID
// @Tags actor
// @Param        id   path      int  true  "Actor ID"
// @Success 200
// @Router /actor/{id} [delete]
// @Security BasicAuth
func (handler ActorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !authWriter(w, r, handler.repo) {
		return
	}
	id, err := getIdFromPath(w, r)
	if err != nil {
		return
	}
	err = handler.repo.ActorRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Получение списка всех актеров
// @Tags actor
// @Produce  json
// @Success 200 {array} model.Actor
// @Router /actors [get]
// @Security BasicAuth
func (handler ActorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if !authReader(w, r, handler.repo) {
		return
	}
	actors, err := handler.repo.ActorRepo.GetAll("name", false)
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
