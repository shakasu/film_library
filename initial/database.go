package initial

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func Db(driver string, source string) (*sql.DB, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS actors (id serial PRIMARY KEY, name varchar(255) not null, gender varchar(1) not null, date_of_birth date not null)",
	); err != nil {
		log.Fatal(err)
	}
	return db, err
}
