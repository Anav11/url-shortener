package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/utils"
)

func SessionMiddleware(conf app.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session")

		if cookie == "" || err != nil {
			encryptedId, err := utils.Encrypt(uuid.New().String(), conf.SecretKey)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.SetCookie("session", encryptedId, 3600, "/", conf.ServerAddress, false, false)
		}

		ctx.Next()
	}
}