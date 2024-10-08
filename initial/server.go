package initial

import (
	"log"
	"net/http"
	"time"
)

func Server(routs *http.ServeMux, err error, addr string) error {
	log.Printf("Документация http://localhost%s/swagger/index.html\n", addr)

	server := &http.Server{
		Addr:           addr,
		Handler:        routs,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	err = server.ListenAndServe()
	return err
}
