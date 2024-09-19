package repository

import (
	"database/sql"
	"film_library/model"
)

var (
	filmTable         = "films"
	filmColumns       = []string{"name", "description", "release_date", "rating", "actors"}
	filmColumnsWithId = []string{"id", "name", "description", "release_date", "rating", "actors"}
)

type FilmRepository struct {
	db *sql.DB
}

func NewFilmRepository(db *sql.DB) *FilmRepository {
	return &FilmRepository{db: db}
}

func (r *FilmRepository) Add(film *model.Film) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *FilmRepository) Update(film *model.Film) error {
	//TODO implement me
	panic("implement me")
}

func (r *FilmRepository) Delete(filmId int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *FilmRepository) GetAll() ([]*model.Film, error) {
	//TODO implement me
	panic("implement me")
}
