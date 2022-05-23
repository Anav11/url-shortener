package server

import (
	"flag"
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
	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		return
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "a localhost:8080")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "b http://localhost:8080")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "f ./urls_db.csv")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "d postgres://username:password@host:port/database")
	flag.Parse()

	s := storage.New(cfg)
	r := router.Router(cfg, s)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		if err := s.Destruct(cfg); err != nil {
			fmt.Printf("ERROR: %s", err)
		}
		os.Exit(0)
	}()

	if err := r.Run(cfg.ServerAddress); err != nil {
		panic(err)
	}
}
