package repository

import (
	"database/sql"
	"film_library/model"
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{db: db}
}

func (r *ActorRepository) Add(actorId int, actor *model.Actor) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) Update(actorId int, actor *model.Actor) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) Delete(ActorId int) (*model.Actor, error) {
	return nil, nil
}
