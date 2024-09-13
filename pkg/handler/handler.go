package handler

import (
	"film_library/pkg/repository"
	"net/http"
)

type Handler struct {
	actorHandler
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{actorHandler: NewActorHandler(repo)}
}

type actorHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

func InitRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	//delete
	mux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("home"))
	})

	//добавить актера (имя пол дата рождения)
	mux.HandleFunc("POST /actor", h.actorHandler.Add)
	//изменить информациою об актере
	mux.HandleFunc("PUT /actor/{id}", h.actorHandler.Update)
	//удалить инф об актере
	mux.HandleFunc("DELETE /actor/{id}", h.actorHandler.Delete)
	//список актеров, для каждого есть список фильмов
	mux.HandleFunc("GET /actors", h.actorHandler.GetAll)

	//добавить инф по фильму
	//mux.HandleFunc("POST /film", h.Put)
	////изм инф по фильму
	//mux.HandleFunc("PUT /film/{id}", h.Put)
	////удл инф по фильму
	//mux.HandleFunc("DELETE /film/{id}", h.Put)
	////список фильмов с возм сортировки по назв рейтингу и дате выпуска (по рейтинг по убыванию дефорт)
	//mux.HandleFunc("GET /film", h.Put)
	//// поиск фильма по фрагменту названия, по фрагменту имени актера
	//mux.HandleFunc("GET /film", h.Put)

	// апи закрыт авторизацией, две роли : одна на чтение, другая на все (соответвие польз и адм задается через бд)
	return mux
}
