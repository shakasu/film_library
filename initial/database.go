package initial

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

const (
	createActors = "CREATE TABLE IF NOT EXISTS actors (id serial PRIMARY KEY, name varchar(255) not null, gender varchar(1) not null, date_of_birth date not null)"
	createFilms  = "CREATE TABLE IF NOT EXISTS films(id serial PRIMARY KEY, name varchar(150) not null, description varchar(1000) not null,release_date date not null, rating smallserial not null)"
	createLink   = "CREATE TABLE IF NOT EXISTS actor_film(id SERIAL PRIMARY KEY, actor_id INTEGER NOT NULL REFERENCES actors, film_id INTEGER NOT NULL REFERENCES films, UNIQUE (actor_id, film_id))"
)

func Db(driver string, source string) (*sql.DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	if _, err := db.Exec(createActors); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(createFilms); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(createLink); err != nil {
		log.Fatal(err)
	}

	return db, err
}
