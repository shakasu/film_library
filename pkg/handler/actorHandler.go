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

func (a ActorHandler) Add(w http.ResponseWriter, r *http.Request) {
	var actorReq model.Actor

	err := decodeJSONBody(w, r, &actorReq)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	id, err := a.repo.Add(&actorReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	actorReq.Id = id
	j, _ := json.Marshal(actorReq)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	if err != nil {
		return
	}

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
