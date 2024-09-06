package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

type config struct {
	Port string `yaml:"port"`
	Db   struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Dbname   string `yaml:"dbname"`
	}
}

func main() {
	var c config
	c.getConfig()

	fmt.Println(c)
	//srv := new(film_library.Server)

	// Регистрируем два новых обработчика и соответствующие URL-шаблоны в
	// маршрутизаторе servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет из Snippetbox"))
	})

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

	//go func() {
	//	if err := srv.Run(c.Port, mux); err != nil {
	//		fmt.Printf("error occured while running http server: %s", err.Error())
	//	}
	//}()
}

func (cfg *config) getConfig() *config {

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
