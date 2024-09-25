package repository

import (
	"database/sql"
	"film_library/model"
)

type Repository struct {
	ActorRepo CrudRepository[model.ActorDto, model.Actor]
	FilmRepo  CrudRepository[model.FilmDto, model.Film]
	AuthRepo  *AuthRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		//ActorRepo: NewActorRepository(db),
		FilmRepo: NewFilmRepository(db),
		AuthRepo: NewAuthRepository(db),
	}
}

type CrudRepository[T any, G any] interface {
	Add(*T) (*G, error)
	Update(*T, int64) (*G, error)
	Delete(int64) error
	GetAll() ([]*G, error)
}
