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

//
//package main
//
//import (
//"database/sql"
//"film_library/pkg/handler"
//"film_library/pkg/repository"
//"fmt"
//_ "github.com/lib/pq"
//"strconv"
//
////"gopkg.in/yaml.v3"
//"log"
//"net/http"
//"os"
//"time"
//)
//
//type db struct {
//	host     string //`yaml:"host"`
//	port     int    //`yaml:"port"`
//	user     string //`yaml:"user"`
//	password string //`yaml:"password"`
//	dbname   string //`yaml:"dbname"`
//	sslmode  string //`yaml:"sslmode"`
//	driver   string //`yaml:"driver"`
//}
//
//type config struct {
//	port int //`yaml:"port"`
//	db   db
//}
//
//func main() {
//	var c config
//	c.initConfig()
//
//	log.Print(c)
//
//	db, err := initDb(c)
//	if err != nil {
//		log.Fatal(err)
//	}
//	routs := handler.NewHandler(repository.NewRepository(db))
//	err = initServer(c, routs, err)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func initServer(c config, routs *handler.Handler, err error) error {
//	addr := fmt.Sprintf(":%d", c.port)
//	log.Printf("Запуск веб-сервера на http://localhost%s", addr)
//
//	server := &http.Server{
//		Addr:           addr,
//		Handler:        handler.InitRoutes(routs),
//		MaxHeaderBytes: 1 << 20,
//		ReadTimeout:    10 * time.Second,
//		WriteTimeout:   10 * time.Second,
//	}
//
//	err = server.ListenAndServe()
//	return err
//}
//
//func initDb(c config) (*sql.DB, error) {
//	dbSource := fmt.Sprintf(
//		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
//		c.db.host, c.db.port, c.db.user, c.db.password, c.db.dbname, c.db.sslmode)
//
//	db, err := sql.Open(c.db.driver, dbSource)
//	if err != nil {
//		log.Fatal("error connecting to the database: ", err)
//	}
//
//	if _, err := db.Exec(
//		"CREATE TABLE IF NOT EXISTS actors (id serial PRIMARY KEY, name varchar(255) not null, gender varchar(1) not null, date_of_birth date not null)",
//	); err != nil {
//		log.Fatal(err)
//	}
//	return db, err
//}
//
//func (cfg *config) initConfig() *config {
//	dbPort, err := strconv.Atoi(os.Getenv("DB.PORT"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	port, err := strconv.Atoi(os.Getenv("PORT"))
//	if err != nil {
//		log.Fatal(err)
//	}
//	return &config{
//		port: port,
//		db: db{
//			host:     os.Getenv("DB.HOST"),
//			port:     dbPort,
//			user:     os.Getenv("DB.USER"),
//			password: os.Getenv("DB.PASSWORD"),
//			dbname:   os.Getenv("DB.DBNAME"),
//			sslmode:  os.Getenv("DB.SSLMODE"),
//			driver:   os.Getenv("DB.DRIVER"),
//		},
//	}
//
//	//yamlFile, err := os.ReadFile("config/config.yml")
//	//if err != nil {
//	//	log.Printf("yamlFile.Get err   #%v ", err)
//	//}
//	//err = yaml.Unmarshal(yamlFile, cfg)
//	//if err != nil {
//	//	log.Fatalf("Unmarshal: %v", err)
//	//}
//	//
//	//return cfg
//}
