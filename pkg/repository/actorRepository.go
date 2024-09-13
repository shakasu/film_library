package repository

import (
	"database/sql"
	"film_library/model"
	"github.com/Masterminds/squirrel"
	"log"
	"time"
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{db: db}
}

func (r *ActorRepository) Add(actor *model.Actor) (int64, error) {

	rawInsert := squirrel.
		Insert("actors").
		Columns("name", "gender", "date_of_birth").
		Values(actor.Name, actor.Gender, actor.DateOfBirth).Suffix("RETURNING id")
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

func (r *ActorRepository) Update(actorId int64, actor *model.Actor) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) Delete(actorId int64) (*model.Actor, error) {
	return nil, nil
}

func (r *ActorRepository) GetAll() ([]*model.Actor, error) {
	selectActors := squirrel.
		Select("id", "name", "gender", "date_of_birth").
		From("actors")

	query, _, err := selectActors.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
