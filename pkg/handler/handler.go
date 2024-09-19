package handler

import (
	"film_library/pkg/repository"
	"net/http"
)

type Handler struct {
	actorHandler crudHandler
	filmHandler  crudHandler
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

func InitRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /actor", h.actorHandler.Add)
	mux.HandleFunc("PUT /actor/{id}", h.actorHandler.Update)
	mux.HandleFunc("DELETE /actor/{id}", h.actorHandler.Delete)
	mux.HandleFunc("GET /actors", h.actorHandler.GetAll)

	//добавить инф по фильму
	mux.HandleFunc("POST /film", h.filmHandler.Add)
	////изм инф по фильму
	mux.HandleFunc("PUT /film/{id}", h.filmHandler.Update)
	////удл инф по фильму
	mux.HandleFunc("DELETE /film/{id}", h.filmHandler.Delete)
	////список фильмов с возм сортировки по назв рейтингу и дате выпуска (по рейтинг по убыванию дефорт)
	mux.HandleFunc("GET /films", h.filmHandler.GetAll)
	//// поиск фильма по фрагменту названия, по фрагменту имени актера
	//mux.HandleFunc("GET /film", h.filmHandler.GetByFilter)

	// апи закрыт авторизацией, две роли : одна на чтение, другая на все (соответвие польз и адм задается через бд)
	return mux
}
