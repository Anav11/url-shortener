package server

import (
	"fmt"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/router"
	"github.com/Anav11/url-shortener/internal/app/storage"
)

func Start(port int) {
	c := app.Config{
		Host: "http://localhost",
		Port: port,
	}
	s := storage.ConstructStorage()
	r := router.Router(c, s)

	r.Run(fmt.Sprintf(":%d", c.Port))
}
