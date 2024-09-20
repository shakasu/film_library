package handler

import (
	"encoding/json"
	"errors"
	"film_library/model"
	"film_library/pkg/repository"
	"film_library/utils"
	"log"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	repo *repository.Repository
}

func NewActorHandler(repo *repository.Repository) *ActorHandler {
	return &ActorHandler{repo: repo}
}

func (handler ActorHandler) Add(w http.ResponseWriter, r *http.Request) {
	var actorReq model.Actor

	err := utils.DecodeJSONBody(w, r, &actorReq)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	id, err := handler.repo.ActorRepo.Add(&actorReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	actorReq.Id = id
	j, err := json.Marshal(actorReq)
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

func (handler ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	var actorReq model.Actor

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.DecodeJSONBody(w, r, &actorReq)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	actorReq.Id = id

	err = handler.repo.ActorRepo.Update(&actorReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	actorReq.Id = id
	j, err := json.Marshal(actorReq)
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

func (handler ActorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = handler.repo.ActorRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler ActorHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	actors, err := handler.repo.ActorRepo.GetAll()
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
