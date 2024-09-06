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
		Password string `yaml:"password"`
	}
}

func main() {
	var c config
	c.getConfig()

	fmt.Println(c)
	//db := film_library.NewDB("postgres", "user=root dbname=tutorial password=root sslmode=disable")
	//Read request body
	//productRequest := Product{}
	//err := helper.ReadFromRequestBody(request, &productRequest)
	//if err != nil {
	//	helper.WriteErrToResponseBody(writer, err)
	//	return
	//}
	//
	////Query to insert data
	//SQL := `INSERT INTO "products" (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	//err = s.db.QueryRow(SQL, productRequest.Name, productRequest.Price, productRequest.Stock).Scan(&productRequest.ID)
	//if err != nil {
	//	helper.WriteErrToResponseBody(writer, err)
	//	return
	//}
	//
	////Write response
	//helper.WriteToResponseBody(writer, productRequest)

	//srv := new(film_library.Server)

	// Регистрируем два новых обработчика и соответствующие URL-шаблоны в
	// маршрутизаторе servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/", funcName())

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

	//go func() {
	//	if err := srv.Run(c.Port, mux); err != nil {
	//		fmt.Printf("error occured while running http server: %s", err.Error())
	//	}
	//}()
}

func funcName() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет из Snippetbox"))
	}
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
