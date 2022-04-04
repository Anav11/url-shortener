package router

import (
	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/handlers"
	"github.com/Anav11/url-shortener/internal/app/storage"
	"github.com/gin-gonic/gin"
)

func Router(c app.Config, s *storage.Storage) *gin.Engine {
	r := gin.Default()
	h := handlers.Handler{
		Config: c,
		Storage: s,
	}
	r.GET("/:ID", h.GetHandler)
	r.POST("/", h.PostHandler)

	return r
}
