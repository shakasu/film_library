package repository

import (
	"database/sql"
	"film_library/model"
	"github.com/Masterminds/squirrel"
	"log"
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{db: db}
}

func (r *ActorRepository) Add(actor *model.Actor) (int64, error) {

	sql := squirrel.
		Insert("actors").
		Columns("name", "gender", "date_of_birth").
		Values(actor.Name, actor.Gender, actor.DateOfBirth).Suffix("RETURNING id")
	query, args, err := sql.PlaceholderFormat(squirrel.Dollar).ToSql()
	log.Println(sql.PlaceholderFormat(squirrel.Dollar).ToSql())
	sql.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ActorRepository) Update(actorId int64, actor *model.Actor) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) Delete(actorId int64) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) GetAll() ([]*model.Actor, error) {
	//TODO implement me
	panic("implement me")
}
