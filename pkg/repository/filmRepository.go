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
	filmTable         = "films"
	filmColumns       = []string{"name", "description", "release_date", "rating"}
	filmColumnsWithId = []string{"id", "name", "description", "release_date", "rating"}
)

type FilmRepo struct {
	db *sql.DB
}

func NewFilmRepository(db *sql.DB) *FilmRepo {
	return &FilmRepo{db: db}
}

func (r *FilmRepo) Add(film *model.FilmDto) (*model.Film, error) {
	for _, actorId := range film.ActorIds {
		exists, err := isRecordExist(actorId, actorTable, r.db)
		if !exists || err != nil {
			return nil, errors.New(fmt.Sprintf("Actor with id = %d not exist", actorId))
		}
	}

	rawFilmInsert := squirrel.
		Insert(filmTable).
		Columns(filmColumns...).
		Values(film.Name, film.Description, film.ReleaseDate, film.Rating).
		Suffix("RETURNING id")
	query, args, err := rawFilmInsert.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	var id int64
	err = r.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return r.handleActors(film, err, id)
}

func (r *FilmRepo) Update(film *model.FilmDto, id int64) (*model.Film, error) {
	for _, actorId := range film.ActorIds {
		exists, err := isRecordExist(actorId, actorTable, r.db)
		if !exists || err != nil {
			return nil, errors.New(fmt.Sprintf("Actor with id = %d not exist", actorId))
		}
	}

	exists, err := isRecordExist(id, filmTable, r.db)
	if !exists || err != nil {
		return nil, errors.New(fmt.Sprintf("Film with id = %d not exist", id))
	}
	err = r.deleteActorsLink(id)
	if err != nil {
		return nil, err
	}
	rawUpdate := squirrel.
		Update(filmTable).
		Set(filmColumns[0], film.Name).
		Set(filmColumns[1], film.Description).
		Set(filmColumns[2], film.ReleaseDate).
		Set(filmColumns[3], film.Rating).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id")
	query, args, err := rawUpdate.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}
	var idUpdated int64
	err = r.db.QueryRow(query, args...).Scan(&idUpdated)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Film with id = %d not exist", id))
	}

	return r.handleActors(film, err, id)
}

func (r *FilmRepo) Delete(id int64) error {
	exists, err := isRecordExist(id, filmTable, r.db)
	if !exists || err != nil {
		return errors.New(fmt.Sprintf("Film with id = %d not exist", id))
	}
	err = r.deleteActorsLink(id)
	if err != nil {
		return err
	}
	rawDelete := squirrel.Delete(filmTable).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id")
	query, args, err := rawDelete.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	var idDeleted int64
	err = r.db.QueryRow(query, args...).Scan(&idDeleted)
	if err != nil {
		return errors.New(fmt.Sprintf("Film with id = %d not exist", id))
	}
	return nil
}

func (r *FilmRepo) GetAll() ([]*model.Film, error) {
	selectFilms := squirrel.
		Select(filmColumnsWithId...).
		From(filmTable)

	query, _, err := selectFilms.ToSql()
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

	var films []*model.Film
	for rows.Next() {
		film := &model.Film{}
		err := rows.Scan(&film.Id, &film.Name, &film.Description, &film.ReleaseDate, &film.Rating)
		if err != nil {
			return nil, err
		}
		releaseDateTime, err := time.Parse(time.RFC3339, film.ReleaseDate)
		if err != nil {
			return nil, err
		}
		film.ReleaseDate = releaseDateTime.Format(time.DateOnly)

		linkedActorIds, err := r.findLinkedActorIds(film.Id)
		if err != nil {
			return nil, err
		}
		linkedActors, err := r.findActorsByIds(linkedActorIds)
		film.Actors = linkedActors

		films = append(films, film)

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return films, nil

}

func (r *FilmRepo) findActorsByIds(actorIds []int64) ([]*model.Actor, error) {
	selectActors := squirrel.
		Select("*").
		From(actorTable).
		Where(squirrel.Eq{"id": actorIds})

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

	var actors []*model.Actor
	for rows.Next() {
		actor := &model.Actor{}
		err := rows.Scan(&actor.Id, &actor.Name, &actor.Gender, &actor.DateOfBirth)
		if err != nil {
			return nil, err
		}
		dateOfBirthTime, err := time.Parse(time.RFC3339, actor.DateOfBirth)
		if err != nil {
			return nil, err
		}
		actor.DateOfBirth = dateOfBirthTime.Format(time.DateOnly)
		actors = append(actors, actor)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return actors, nil
}

func (r *FilmRepo) deleteActorsLink(id int64) error {
	exists, _ := isRecordExist(id, "actor_film", r.db)
	if exists {
		rawLinkDelete := squirrel.
			Delete("actor_film").
			Where(squirrel.Eq{"film_id": id}).
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

func (r *FilmRepo) findLinkedActorIds(filmId int64) ([]int64, error) {
	rawLinkSelect := squirrel.
		Select("actor_id").
		From("actor_film").
		Where(squirrel.Eq{"film_id": filmId})

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

func (r *FilmRepo) handleActors(film *model.FilmDto, err error, id int64) (*model.Film, error) {
	var actors []*model.Actor
	if film.ActorIds != nil {
		actors, err = r.findActorsByIds(film.ActorIds)
		if err != nil {
			return nil, err
		}

		rawLinkInsert := squirrel.
			Insert("actor_film").
			Columns("actor_id", "film_id")
		for _, actorId := range film.ActorIds {
			rawLinkInsert = rawLinkInsert.Values(actorId, id)
		}

		query, args, err := rawLinkInsert.PlaceholderFormat(squirrel.Dollar).ToSql()
		if err != nil {
			return nil, err
		}

		err = r.db.QueryRow(query, args...).Err()
		if err != nil {
			return nil, err
		}
	}

	return &model.Film{
		Id:          id,
		Name:        film.Name,
		Description: film.Description,
		ReleaseDate: film.ReleaseDate,
		Rating:      film.Rating,
		Actors:      actors,
	}, nil
}
