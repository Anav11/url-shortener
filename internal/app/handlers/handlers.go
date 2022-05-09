package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/storage"
	"github.com/Anav11/url-shortener/internal/app/utils"
)

type Handler struct {
	Config  app.Config
	Storage storage.Repository
}

type ShortenerResponseJSON struct {
	Result string `json:"result"`
}

type ShortenerRequestJSON struct {
	URL string `json:"url"`
}

type UserURLsJSON struct {
	ShortURL	string `json:"short_url"`
	OriginalURL	string `json:"original_url"`
}

func (h Handler) GetHandler(ctx *gin.Context) {
	ID := ctx.Param("ID")
	if ID == "" {
		ctx.String(http.StatusBadRequest, "")
		return
	}

	initialURL, err := h.Storage.GetURL(ID)
	if err != nil {
		ctx.String(http.StatusNotFound, "")
		return
	}

	ctx.Header("Content-Type", "text/plain")
	ctx.Redirect(http.StatusTemporaryRedirect, initialURL)
}

func (h Handler) PostHandler(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	ID, err := createURL(h, ctx, string(body))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, ID)

	ctx.Header("Content-Type", "text/plain")
	ctx.String(http.StatusCreated, "%s", shortURL)
}

func (h Handler) PostHandlerJSON(ctx *gin.Context) {
	var req ShortenerRequestJSON
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	ID, err := createURL(h, ctx, req.URL)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, ID)
	res := ShortenerResponseJSON{Result: shortURL}

	ctx.JSON(http.StatusCreated, res)
}

func (h Handler) GetUserURLsHandler(ctx *gin.Context) {
	userID, err := ctx.Cookie("session")
	if err != nil {
		ctx.String(http.StatusUnprocessableEntity, "cookies were not set")
		return
	}

	userDecryptID, err := utils.Decrypt(userID, h.Config.SecretKey)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	userShortURLs := h.Storage.GetUserShortURLs(userDecryptID)
	if len(userShortURLs) == 0 {
		ctx.JSON(http.StatusNoContent, "{}")
		return
	}

	var userURLsJSON []UserURLsJSON
	for _, userShorted := range userShortURLs {
		shortURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, userShorted.ID)
		userURLsJSON = append(userURLsJSON, UserURLsJSON{shortURL, userShorted.OriginalURL})
	}

	ctx.JSON(http.StatusOK, userURLsJSON)
}

func (h Handler) PingDBHandler(ctx *gin.Context) {
	err := h.Storage.Ping()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "database is running")
}

func createURL(h Handler, ctx *gin.Context, URL string) (shortURLID string, error error) {
	userEncryptID, err := ctx.Cookie("session")
	shortURLID = uuid.New().String()

	if userEncryptID != "" && err == nil {
		userDecryptID, err := utils.Decrypt(userEncryptID, h.Config.SecretKey)
		if err != nil {
			return "", err
		}

		if err := h.Storage.AddURL(shortURLID, URL, userDecryptID); err != nil {
			ctx.String(http.StatusBadRequest, "")
			return
		}
	} else {
		if err := h.Storage.AddURL(shortURLID, URL, ""); err != nil {
			ctx.String(http.StatusBadRequest, "")
			return
		}
	}

	return shortURLID, nil
}
