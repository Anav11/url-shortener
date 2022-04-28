package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/handlers"
	"github.com/Anav11/url-shortener/internal/app/storage"
)

func Router(c app.Config, s *storage.Storage) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	h := handlers.Handler{
		Config: c,
		Storage: s,
	}
	r.GET("/:ID", h.GetHandler)
	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.PostHandlerJSON)

	return r
}
