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

type FilmRepository struct {
	db *sql.DB
}

func NewFilmRepository(db *sql.DB) *FilmRepository {
	return &FilmRepository{db: db}
}

func (r *FilmRepository) Add(film *model.Film) (int64, error) {
	rawInsert := squirrel.
		Insert(filmTable).
		Columns(filmColumns...).
		Values(film.Name, film.Description, film.ReleaseDate, film.Rating).
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

func (r *FilmRepository) Update(film *model.Film) error {
	rawUpdate := squirrel.
		Update(filmTable).
		Set(filmColumns[0], film.Name).
		Set(filmColumns[1], film.Description).
		Set(filmColumns[2], film.ReleaseDate).
		Set(filmColumns[3], film.Rating).
		Where(squirrel.Eq{"id": film.Id}).
		Suffix("RETURNING id")
	query, args, err := rawUpdate.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	var idDeleted int64
	err = r.db.QueryRow(query, args...).Scan(&idDeleted)
	if err != nil {
		return errors.New(fmt.Sprintf("Film with id = %d not exist", film.Id))
	}
	return nil
}

func (r *FilmRepository) Delete(filmId int64) error {
	rawDelete := squirrel.Delete(filmTable).
		Where(squirrel.Eq{"id": filmId}).
		Suffix("RETURNING id")
	query, args, err := rawDelete.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	var idDeleted int64
	err = r.db.QueryRow(query, args...).Scan(&idDeleted)
	if err != nil {
		return errors.New(fmt.Sprintf("Actor with id = %d not exist", filmId))
	}
	return nil
}

func (r *FilmRepository) GetAll() ([]*model.Film, error) {
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
			log.Fatal(err)
		}
		releaseDateTime, err := time.Parse(time.RFC3339, film.ReleaseDate)
		if err != nil {
			log.Fatal(err)
		}
		film.ReleaseDate = releaseDateTime.Format(time.DateOnly)
		films = append(films, film)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return films, nil
}
