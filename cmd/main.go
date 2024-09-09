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
		Dbname   string `yaml:"name"`
		Sslmode  string `yaml:"sslmode"`
		Driver   string `yaml:"driver"`
	}
}

func main() {
	var c config
	c.initConfig()
	log.Print(c)
	dbSource := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Db.Host, c.Db.Port, c.Db.User, c.Db.Password, c.Db.Dbname, c.Db.Sslmode)
	db, err := sql.Open(c.Db.Driver, dbSource)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	routs := handler.NewHandler(repository.NewRepository(db))
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
	log.Fatal(err)
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
