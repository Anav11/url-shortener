package middleware

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/utils"
)

func SessionMiddleware(conf app.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session")

		if cookie == "" || err != nil {
			encryptedID, err := utils.Encrypt(uuid.New().String(), conf.SecretKey)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			ctx.Request.AddCookie(&http.Cookie{
				Name:     "session",
				Value:    url.QueryEscape(encryptedID),
			})

			ctx.SetCookie("session", encryptedID, 3600, "/", conf.ServerAddress, false, false)
		}

		ctx.Next()
	}
}