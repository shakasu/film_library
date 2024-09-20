package initial

import (
	"film_library/pkg/handler"
	"log"
	"net/http"
	"time"
)

func Server(routs *handler.Handler, err error, addr string) error {
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
