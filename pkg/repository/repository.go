package repository

import (
	"database/sql"
	"film_library/model"
)

type Repository struct {
	ActorRepo CrudRepository[model.Actor]
	FilmRepo  CrudRepository[model.Film]
	AuthRepo  *AuthRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ActorRepo: NewActorRepository(db),
		FilmRepo:  NewFilmRepository(db),
		AuthRepo:  NewAuthRepository(db),
	}
}

type CrudRepository[T any] interface {
	Add(*T) (int64, error)
	Update(*T) error
	Delete(int64) error
	GetAll() ([]*T, error)
}
