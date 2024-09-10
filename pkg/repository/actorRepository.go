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

func (r *ActorRepository) Add(actor *model.Actor) (int, error) {
	log.Print(actor)
	sql, args, err := squirrel.
		Insert("actor_repository").
		Columns("name", "gender", "date_of_birth").
		Values(actor.Name, actor.Gender, actor.DateOfBirth).
		ToSql()
	if err != nil {
		return 0, err
	}
	log.Println("sql" + sql)
	log.Println(args)
	//sql, args, err := sq.
	//    Insert("users").Columns("name", "age").
	//    Values("moe", 13).Values("larry", sq.Expr("? + 5", 12)).
	//    ToSql()
	//SQL := "INSERT INTO actors (name, gender, date_of_birth) VALUES ($1, $2, $3) RETURNING id"
	err = r.db.QueryRow(sql, actor.Name, actor.Gender, actor.DateOfBirth).Scan(&actor.Id)
	if err != nil {
		return 0, err
	}
	//row := r.db.QueryRow(SQL, actor.Name, actor.Gender, actor.DateOfBirth)
	//if row.Err() != nil {
	//	return 0, err
	//}
	return actor.Id, nil
}

func (r *ActorRepository) Update(actorId int, actor *model.Actor) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) Delete(ActorId int) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) GetAll() ([]*model.Actor, error) {
	//TODO implement me
	panic("implement me")
}
