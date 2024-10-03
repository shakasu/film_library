package handler

import (
	"encoding/json"
	"errors"
	"film_library/pkg/repository"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type malformedRequest struct {
	Status int
	Msg    string
}

type Handler struct {
	actorHandler crudHandler
	filmHandler  crudAndSearchHandler
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		actorHandler: NewActorHandler(repo),
		filmHandler:  NewFilmHandler(repo),
	}
}

type crudHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type searchHandler interface {
	searchBy(w http.ResponseWriter, r *http.Request)
}

type crudAndSearchHandler interface {
	crudHandler
	searchHandler
}

func InitRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /actor", h.actorHandler.Add)
	mux.HandleFunc("PUT /actor/{id}", h.actorHandler.Update)
	mux.HandleFunc("DELETE /actor/{id}", h.actorHandler.Delete)
	mux.HandleFunc("GET /actors", h.actorHandler.GetAll)

	mux.HandleFunc("POST /film", h.filmHandler.Add)
	mux.HandleFunc("PUT /film/{id}", h.filmHandler.Update)
	mux.HandleFunc("DELETE /film/{id}", h.filmHandler.Delete)
	mux.HandleFunc("GET /films", h.filmHandler.GetAll)

	mux.HandleFunc("GET /film/search/{fragment}", h.filmHandler.searchBy)
	return mux
}

func (mr *malformedRequest) Error() string {
	return mr.Msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{Status: http.StatusUnsupportedMediaType, Msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{Status: http.StatusRequestEntityTooLarge, Msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{Status: http.StatusBadRequest, Msg: msg}
	}

	return nil
}

func authWriter(w http.ResponseWriter, r *http.Request, repo *repository.Repository) bool {
	username, password, ok := r.BasicAuth()
	if ok {
		auth, role := repo.AuthRepo.Authorize(username, password)
		if !auth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return false
		}
		if role != "ADMIN" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return false
		}
	}
	return true
}

func authReader(w http.ResponseWriter, r *http.Request, repo *repository.Repository) bool {
	username, password, ok := r.BasicAuth()
	if ok {
		auth, _ := repo.AuthRepo.Authorize(username, password)
		if !auth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return false
		}
	}
	return true
}

func getIdFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0, nil
	}
	return id, err
}
