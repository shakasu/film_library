package handler

import (
	"film_library/pkg/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func InitRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	//добавить актера (имя пол дата рождения)
	mux.HandleFunc("POST /api/actor", h.services.Add())
	//изменить информациою об актере
	mux.HandleFunc("PUT /api/actor/{id}", h.Get)
	//удалить инф об актере
	mux.HandleFunc("DELETE /api/actor/{id}", h.services.actorService.Id)
	//добавить инф по фильму
	//mux.HandleFunc("POST /api/film", h.Put)
	////изм инф по фильму
	//mux.HandleFunc("PUT /api/film/{id}", h.Put)
	////удл инф по фильму
	//mux.HandleFunc("DELETE /api/film/{id}", h.Put)
	////список фильмов с возм сортировки по назв рейтингу и дате выпуска (по рейтинг по убыванию дефорт)
	//mux.HandleFunc("GET /api/film", h.Put)
	//// поиск фильма по фрагменту названия, по фрагменту имени актера
	//mux.HandleFunc("GET /api/film", h.Put)
	////список актеров, для каждого есть список фильмов
	//mux.HandleFunc("GET /api/actor", h.Put)

	// апи закрыт авторизацией, две роли : одна на чтение, другая на все (соответвие польз и адм задается через бд)
	return mux
}
