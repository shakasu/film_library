package main

import (
	"database/sql"
	"film_library/pkg/handler"
	"film_library/pkg/repository"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	Port string `yaml:"port"`
	Db   struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
		Driver   string `yaml:"driver"`
	}
}

func main() {
	var c config
	c.initConfig()

	db, err := initDb(c)
	if err != nil {
		log.Fatal(err)
	}
	routs := handler.NewHandler(repository.NewRepository(db))
	err = initServer(c, routs, err)
	if err != nil {
		log.Fatal(err)
	}
}

func initServer(c config, routs *handler.Handler, err error) error {
	addr := ":" + c.Port
	log.Printf("Запуск веб-сервера на http://localhost%s", addr)

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

func initDb(c config) (*sql.DB, error) {
	dbSource := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Db.Host, c.Db.Port, c.Db.User, c.Db.Password, c.Db.Dbname, c.Db.Sslmode)

	db, err := sql.Open(c.Db.Driver, dbSource)
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

func (cfg *config) initConfig() *config {

	yamlFile, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return cfg
}
