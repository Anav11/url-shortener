package handlers

import (
	"encoding/json"
	"errors"
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

func (h Handler) GetHandler(ctx *gin.Context) {
	ID := ctx.Param("ID")
	if ID == "" {
		ctx.String(http.StatusBadRequest, "empty ID")
		return
	}

	initialURL, err := h.Storage.GetURL(ID)
	if err != nil {
		var due *storage.DeletedURLError
		if errors.As(err, &due) {
			ctx.String(http.StatusGone, "")
		}

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

	URL := string(body)
	ID, err := createURL(h, ctx, URL)
	if err != nil {
		var ude *storage.URLDuplicateError
		if errors.As(err, &ude) {
			existID, _ := h.Storage.GetShortByOriginal(URL)
			existURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, existID)
			ctx.String(http.StatusConflict, existURL)
			return
		}

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
		var ude *storage.URLDuplicateError
		if errors.As(err, &ude) {
			existID, _ := h.Storage.GetShortByOriginal(req.URL)
			existURL := fmt.Sprintf("%s/%s", h.Config.BaseURL, existID)
			res := ShortenerResponseJSON{Result: existURL}

			ctx.JSON(http.StatusConflict, res)
			return
		}

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

func (h Handler) PostBatchHandler(ctx *gin.Context) {
	var batchURLs BatchURLs
	userID, _ := ctx.Cookie("session")

	if err := json.NewDecoder(ctx.Request.Body).Decode(&batchURLs); err != nil {
		ctx.String(http.StatusInternalServerError, "batch decoding error")
		return
	}

	shortURLs := make([]storage.UserShortURL, 0)

	for _, bu := range batchURLs {
		shortURLs = append(shortURLs, storage.UserShortURL{ID: bu.CorrelationID, OriginalURL: bu.OriginalURL, UserID: userID})
	}

	if err := h.Storage.AddBatchURL(shortURLs); err != nil {
		ctx.String(http.StatusInternalServerError, "batch add error")
		return
	}

	batchShortURL := make([]BatchShortURL, 0)

	for _, su := range shortURLs {
		batchShortURL = append(batchShortURL, BatchShortURL{
			CorrelationID: su.ID,
			ShortURL: fmt.Sprintf("%s/%s", h.Config.BaseURL, su.ID),
		})
	}

	ctx.JSON(http.StatusCreated, batchShortURL)
}

func (h Handler) DeleteUserURLsHandler(ctx *gin.Context) {
	var IDs []string
	if err := json.NewDecoder(ctx.Request.Body).Decode(&IDs); err != nil {
		return
	}

	userID, _ := ctx.Cookie("session")
	userDecryptID, err := utils.Decrypt(userID, h.Config.SecretKey)
	if err != nil {
		return
	}

	go func() {
		h.Storage.DeleteUserURLs(IDs, userDecryptID)
	}()

	ctx.String(http.StatusAccepted, "")
}

func createURL(h Handler, ctx *gin.Context, URL string) (shortURLID string, error error) {
	userEncryptID, err := ctx.Cookie("session")
	shortURLID = uuid.New().String()

	if userEncryptID != "" && err == nil {
		userDecryptID, err := utils.Decrypt(userEncryptID, h.Config.SecretKey)
		if err != nil {
			return "", err
		}

		if err := h.Storage.AddURL(storage.UserShortURL{ID: shortURLID, OriginalURL: URL, UserID: userDecryptID}); err != nil {
			return "", err
		}
	} else {
		if err := h.Storage.AddURL(storage.UserShortURL{ID: shortURLID, OriginalURL: URL, UserID: ""}); err != nil {
			return "", err
		}
	}

	return shortURLID, nil
}
