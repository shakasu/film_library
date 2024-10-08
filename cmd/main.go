package main

import (
	_ "film_library/docs"
	"film_library/initial"
	"film_library/pkg/handler"
	"film_library/pkg/repository"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /v2
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
