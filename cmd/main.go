package main

import (
	_ "film_library/docs"
	"film_library/initial"
	"film_library/pkg/handler"
	"film_library/pkg/repository"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
)

// @title Film library golang test task
// @version 1.0
// @description Это приложение управления базой данных "Фильмотека".
// @contact.name   API Support
// @contact.url https://t.me/shakasu
// @host localhost:8080
// @securityDefinitions.basic  BasicAuth
func main() {
	config, err := initial.Cfg()

	db, err := initial.Db(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal(err)
	}
	routs := handler.InitRoutes(
		handler.NewHandler(
			repository.NewRepository(db)))

	routs.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	err = initial.Server(routs, err, config.Port)
	if err != nil {
		log.Fatal(err)
	}
}
