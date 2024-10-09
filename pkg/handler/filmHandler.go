package handler

import (
	"encoding/json"
	"errors"
	"film_library/model"
	"film_library/pkg/repository"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type FilmHandler struct {
	repo *repository.Repository
}

func NewFilmHandler(repo *repository.Repository) *FilmHandler {
	return &FilmHandler{repo: repo}
}

// @Summary Создание фильма с привязкой актеров по ID
// @Tags film
// @Param data body model.FilmDto true "The input film struct"
// @accept  json
// @Produce  json
// @Success 200
// @Router /film [post]
// @Security BasicAuth
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

// @Summary Редактирование фильма по ID с привязкой актеров по ID
// @Tags film
// @Param        id   path      int  true  "Film ID"
// @Param data body model.FilmDto true "The input film struct"
// @accept  json
// @Produce  json
// @Success 200
// @Router /film/{id} [put]
// @Security BasicAuth
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

// @Summary Удаление фильма по ID
// @Tags film
// @Param        id   path      int  true  "Film ID"
// @Success 200
// @Router /film/{id} [delete]
// @Security BasicAuth
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

// @Summary Получение списка всех фильмов
// @Tags film
// @Produce json
// @Success 200 {array} model.Film
// @Param sortBy    query string  false "доступные значения" Enums(name, release_date, rating)
// @Param ascending query boolean false "доступные значения" Enums(true, false)
// @Router /films [get]
// @Security BasicAuth
func (handler FilmHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if !authReader(w, r, handler.repo) {
		return
	}

	sortBy, ascending, err := readGetAllQueryParams(w, r)
	if err != nil {
		return
	}

	films, err := handler.repo.FilmRepo.GetAll(sortBy, ascending)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(films)
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

// @Summary Поиск фильма по фрагменту имени фильма\актера
// @Tags film
// @Param        fragment   path      string  true  "Film ID"
// @Produce json
// @Success 200 {array} model.Film
// @Router /film/search/{fragment} [get]
// @Security BasicAuth
func (handler FilmHandler) searchBy(w http.ResponseWriter, r *http.Request) {
	if !authReader(w, r, handler.repo) {
		return
	}

	fragment := r.PathValue("fragment")
	if fragment == "" {
		msg := "field [searchBy] must not be empty"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	films, err := handler.repo.FilmRepo.SearchBy(fragment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(films)
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

func readGetAllQueryParams(w http.ResponseWriter, r *http.Request) (string, bool, error) {
	var sortBy = r.URL.Query().Get("sortBy")
	var ascendingStr = r.URL.Query().Get("ascending")

	if sortBy == "" && ascendingStr == "" {
		return "rating", false, nil
	}

	if sortBy != "name" && sortBy != "release_date" && sortBy != "rating" {
		msg := fmt.Sprintf(
			"field [sortBy] = [%s] the field can only take 3 values: [name], [release_date], [rating]",
			sortBy)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return "", false, errors.New(msg)
	}
	ascending, err := strconv.ParseBool(ascendingStr)
	if err != nil {
		msg := fmt.Sprintf(
			"field [ascending] = [%t] the field can only take 2 values: [true], [false]",
			ascending)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return "", false, errors.New(msg)
	}

	return sortBy, ascending, err
}
