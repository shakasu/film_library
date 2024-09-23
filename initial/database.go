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

	return db, err
}
