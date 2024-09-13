package repository

import (
	"database/sql"
	"film_library/model"
)

type Repository struct {
	Actor
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{Actor: NewActorRepository(db)}
}

type Actor interface {
	Add(actor *model.Actor) (int64, error)
	Update(actor *model.Actor) error
	Delete(actorId int64) error
	GetAll() ([]*model.Actor, error)
}
