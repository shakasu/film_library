package main

import (
	"database/sql"
	"film_library/pkg/handler"
	"film_library/pkg/repository"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port     string
	dbDriver string
	dbSource string
}

func main() {
	cfg := initConfig()
	log.Println(cfg)
	db, err := initDb(cfg.dbDriver, cfg.dbSource)
	if err != nil {
		log.Fatal(err)
	}
	routs := handler.NewHandler(repository.NewRepository(db))
	err = initServer(routs, err, cfg.port)
	if err != nil {
		log.Fatal(err)
	}
}

func initConfig() config {
	return config{
		port:     ":" + os.Getenv("PORT"),
		dbDriver: os.Getenv("DB_DRIVER"),
		dbSource: os.Getenv("DB_SOURCE"),
	}
}

func initDb(driver string, source string) (*sql.DB, error) {
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

func initServer(routs *handler.Handler, err error, addr string) error {
	log.Printf("Запуск веб-сервера на http://localhost%s\n", addr)

	server := &http.Server{
		Addr:           addr,
		Handler:        handler.InitRoutes(routs),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	err = server.ListenAndServe()
	return err
}
