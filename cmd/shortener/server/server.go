package server

import (
	"github.com/Anav11/url-shortener/internal/app/handlers"
	"log"
	"net/http"
)

func Start(port string) {
	http.HandleFunc("/", handlers.MainHandler)

	log.Println("Server started on port", port)

	http.ListenAndServe(port, nil)
}
