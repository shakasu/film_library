package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Authorize(user string, password string) (bool, string) {
	rawSelect := squirrel.
		Select("(password = crypt($1, password)) AS match", "user_role").
		From("accounts").
		Where("login = $2")
	query, _, _ := rawSelect.PlaceholderFormat(squirrel.Dollar).ToSql()

	var match bool
	var userRole string
	_ = r.db.QueryRow(query, password, user).Scan(&match, &userRole)
	return match, userRole
}
