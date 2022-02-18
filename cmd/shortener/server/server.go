package server

import (
	"log"
	"net/http"
)

func Start(port string) {
	server := &http.Server{
		Addr: "localhost:" + port,
	}
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/{id}", mainHandler)

	log.Println("Server started on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
