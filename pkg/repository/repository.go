package repository

import (
	"database/sql"
	"film_library/model"
	"github.com/Masterminds/squirrel"
)

type Repository struct {
	ActorRepo crudRepository[model.ActorDto, model.Actor]
	FilmRepo  crudAndSearchRepository[model.FilmDto, model.Film]
	AuthRepo  *AuthRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		ActorRepo: NewActorRepository(db),
		FilmRepo:  NewFilmRepository(db),
		AuthRepo:  NewAuthRepository(db),
	}
}

type crudRepository[T any, G any] interface {
	Add(*T) (*G, error)
	Update(*T, int64) (*G, error)
	Delete(int64) error
	GetAll(sortBy string, ascending bool) ([]*G, error)
}

type searchRepository[T any] interface {
	SearchBy(fragment string) ([]*T, error)
}
type crudAndSearchRepository[T any, G any] interface {
	crudRepository[T, G]
	searchRepository[G]
}

func isRecordExist(id int64, table string, db *sql.DB) (bool, error) {
	rawSelect := squirrel.
		Select("1").
		Prefix("SELECT EXISTS (").
		From(table).
		Where(squirrel.Eq{"id": id}).
		Suffix(")")
	query, args, err := rawSelect.PlaceholderFormat(squirrel.Dollar).ToSql()
	var exists bool
	err = db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}
