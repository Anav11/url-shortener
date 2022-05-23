package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/handlers"
	"github.com/Anav11/url-shortener/internal/app/middleware"
	"github.com/Anav11/url-shortener/internal/app/storage"
)

func Router(c app.Config, s storage.Repository) *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.Use(middleware.SessionMiddleware(c))
	h := handlers.Handler{
		Config: c,
		Storage: s,
	}
	r.GET("/:ID", h.GetHandler)
	r.GET("/api/user/urls", h.GetUserURLsHandler)
	r.GET("/ping", h.PingDBHandler)

	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.PostHandlerJSON)
	r.POST("/api/shorten/batch", h.PostBatchHandler)

	return r
}
