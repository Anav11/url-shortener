package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	s := storage.ConstructStorage(c.FileStoragePath)
	r := router.Router(c, s)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		if err := storage.DestructStorage(c.FileStoragePath, s); err != nil {
			fmt.Printf("ERROR: %s", err)
		}
		os.Exit(0)
	}()

	if err := r.Run(c.ServerAddress); err != nil {
		panic(err)
	}
}
