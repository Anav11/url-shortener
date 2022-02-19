package server

import (
	"log"
	"net/http"
)

func Start(port string) {
	http.HandleFunc("/", mainHandler)

	log.Println("Server started on port", port)

	http.ListenAndServe(port, nil)
}
