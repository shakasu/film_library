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

type ActorRepo struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepo {
	return &ActorRepo{db: db}
}

func (r *ActorRepo) Add(actor *model.ActorDto) (*model.Actor, error) {
	rawInsert := squirrel.
		Insert(actorTable).
		Columns(actorColumns...).
		Values(actor.Name, actor.Gender, actor.DateOfBirth).
		Suffix("RETURNING id")
	query, args, err := rawInsert.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	var id int64
	err = r.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &model.Actor{
		Id:          id,
		Name:        actor.Name,
		Gender:      actor.Gender,
		DateOfBirth: actor.DateOfBirth,
	}, nil
}

func (r *ActorRepo) Update(actor *model.ActorDto, id int64) (*model.Actor, error) {
	rawUpdate := squirrel.
		Update(actorTable).
		Set(actorColumns[0], actor.Name).
		Set(actorColumns[1], actor.Gender).
		Set(actorColumns[2], actor.DateOfBirth).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id")
	query, args, err := rawUpdate.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	var idDeleted int64
	err = r.db.QueryRow(query, args...).Scan(&idDeleted)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Actor with id = %d not exist", id))
	}
	return &model.Actor{
		Id:          id,
		Name:        actor.Name,
		Gender:      actor.Gender,
		DateOfBirth: actor.DateOfBirth,
	}, nil
}

func (r *ActorRepo) Delete(id int64) error {
	exists, err := r.isActorExist(id)
	if !exists || err != nil {
		return errors.New(fmt.Sprintf("Actor with id = %d not exist", id))
	}
	err = r.deleteFilmsLink(id)
	if err != nil {
		return err
	}
	rawDelete := squirrel.Delete(actorTable).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id")
	query, args, err := rawDelete.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	var idDeleted int64
	err = r.db.QueryRow(query, args...).Scan(&idDeleted)
	if err != nil {
		return errors.New(fmt.Sprintf("Actor with id = %d not exist", id))
	}
	return nil
}

func (r *ActorRepo) GetAll(_ string, _ bool) ([]*model.Actor, error) {
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

		linkedFilmsIds, err := r.findLinkedFilmIds(actor.Id)
		if err != nil {
			return nil, err
		}
		linkedFilms, err := r.findFilmsByIds(linkedFilmsIds)
		actor.Films = linkedFilms

		actors = append(actors, actor)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return actors, nil
}

func (r *ActorRepo) isActorExist(id int64) (bool, error) {
	rawSelect := squirrel.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("actors").
		Where(squirrel.Eq{"id": id}).
		Suffix(")")
	query, args, err := rawSelect.PlaceholderFormat(squirrel.Dollar).ToSql()
	var exists bool
	err = r.db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}

func (r *ActorRepo) deleteFilmsLink(id int64) error {
	if r.isLinkExist(id) {
		rawLinkDelete := squirrel.
			Delete("actor_film").
			Where(squirrel.Eq{"actor_id": id}).
			Suffix("RETURNING id")

		query, args, err := rawLinkDelete.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return err
		}

		var idDeleted int64
		err = r.db.QueryRow(query, args...).Scan(&idDeleted)
		if err != nil {
			return errors.New(fmt.Sprintf("Actor-film link with id = %d not exist", idDeleted))
		}
	}

	return nil
}

func (r *ActorRepo) isLinkExist(id int64) bool {
	rawSelect := squirrel.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("actor_film").
		Where(squirrel.Eq{"actor_id": id}).
		Suffix(")")
	query, args, _ := rawSelect.PlaceholderFormat(squirrel.Dollar).ToSql()
	var exists bool
	_ = r.db.QueryRow(query, args...).Scan(&exists)
	return exists
}

func (r *ActorRepo) findFilmsByIds(filmIds []int64) ([]*model.Film, error) {
	selectActors := squirrel.
		Select("*").
		From(filmTable).
		Where(squirrel.Eq{"id": filmIds})

	query, args, err := selectActors.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var films []*model.Film
	for rows.Next() {
		film := &model.Film{}
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return nil, err
		}
		releaseDate, err := time.Parse(time.RFC3339, film.ReleaseDate)
		if err != nil {
			return nil, err
		}
		film.ReleaseDate = releaseDate.Format(time.DateOnly)
		films = append(films, film)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return films, nil
}

func (r *ActorRepo) findLinkedFilmIds(actorId int64) ([]int64, error) {
	rawLinkSelect := squirrel.
		Select("film_id").
		From("actor_film").
		Where(squirrel.Eq{"actor_id": actorId})

	query, args, err := rawLinkSelect.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var ids []int64
	rows, err := r.db.Query(query, args...)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
