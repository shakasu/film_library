package main

import (
	"film_library/initial"
	"film_library/pkg/handler"
	"film_library/pkg/repository"
	"log"
)

func main() {
	config, err := initial.Cfg()

	db, err := initial.Db(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal(err)
	}
	routs := handler.NewHandler(repository.NewRepository(db))
	err = initial.Server(routs, err, config.Port)
	if err != nil {
		log.Fatal(err)
	}
}
