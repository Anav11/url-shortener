package server

import (
	"github.com/caarlos0/env/v6"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/router"
	"github.com/Anav11/url-shortener/internal/app/storage"
)

func Start() {
	c := app.Config{}
	if err := env.Parse(&c); err != nil {
		return
	}

	s := storage.ConstructStorage()
	r := router.Router(c, s)

	if err := r.Run(c.ServerAddress); err != nil {
		panic(err)
	}
}
