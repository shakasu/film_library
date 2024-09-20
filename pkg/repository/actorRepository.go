package repository

import (
	"database/sql"
	"errors"
	"film_library/model"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
	"time"
)

var (
	actorTable         = "actors"
	actorColumns       = []string{"name", "gender", "date_of_birth"}
	actorColumnsWithId = []string{"id", "name", "gender", "date_of_birth"}
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{db: db}
}

func (r *ActorRepository) Add(actor *model.Actor) (int64, error) {
	rawInsert := squirrel.
		Insert(actorTable).
		Columns(actorColumns...).
		Values(actor.Name, actor.Gender, actor.DateOfBirth).
		Suffix("RETURNING id")
	query, args, err := rawInsert.PlaceholderFormat(squirrel.Dollar).ToSql()
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

func (r *ActorRepository) Update(actor *model.Actor) error {
	rawUpdate := squirrel.
		Update(actorTable).
		Set(actorColumns[0], actor.Name).
		Set(actorColumns[1], actor.Gender).
		Set(actorColumns[2], actor.DateOfBirth).
		Where(squirrel.Eq{"id": actor.Id}).
		Suffix("RETURNING id")
	query, args, err := rawUpdate.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	var idDeleted int64
	err = r.db.QueryRow(query, args...).Scan(&idDeleted)
	if err != nil {
		return errors.New(fmt.Sprintf("Actor with id = %d not exist", actor.Id))
	}
	return nil
}

func (r *ActorRepository) Delete(actorId int64) error {
	rawDelete := squirrel.Delete(actorTable).
		Where(squirrel.Eq{"id": actorId}).
		Suffix("RETURNING id")
	query, args, err := rawDelete.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	var idDeleted int64
	err = r.db.QueryRow(query, args...).Scan(&idDeleted)
	if err != nil {
		return errors.New(fmt.Sprintf("Actor with id = %d not exist", actorId))
	}
	return nil
}

func (r *ActorRepository) GetAll() ([]*model.Actor, error) {
	selectActors := squirrel.
		Select(actorColumnsWithId...).
		From(actorTable)

	query, _, err := selectActors.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var actors []*model.Actor
	for rows.Next() {
		actor := &model.Actor{}
		err := rows.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.DateOfBirth)
		if err != nil {
			log.Fatal(err)
		}
		dateOfBirthTime, err := time.Parse(time.RFC3339, actor.DateOfBirth)
		if err != nil {
			log.Fatal(err)
		}
		actor.DateOfBirth = dateOfBirthTime.Format(time.DateOnly)
		actors = append(actors, actor)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return actors, nil
}
