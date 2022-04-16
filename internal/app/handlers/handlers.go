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

func (h Handler) GetHandler(ctx *gin.Context) {
	ID := ctx.Param("ID")
	if ID == "" {
		ctx.String(http.StatusBadRequest, "")
		return
	}

	initialURL, err := h.Storage.Get(ID)
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

	ID := uuid.New().String()
	h.Storage.Add(ID, string(body))

	shortURL := fmt.Sprintf("%s:%d/%s", h.Config.Host, h.Config.Port, ID)

	ctx.Header("Content-Type", "text/plain")
	ctx.String(http.StatusCreated, "%s", shortURL)
}

func (h Handler) PostHandlerJSON(ctx *gin.Context) {
	var req ShortenerRequestJSON
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	ID := uuid.New().String()
	if err := h.Storage.Add(ID, req.URL); err != nil {
		ctx.String(http.StatusBadRequest, "")
		return
	}

	shortURL := fmt.Sprintf("%s:%d/%s", h.Config.Host, h.Config.Port, ID)
	res := ShortenerResponseJSON{Result: shortURL}

	ctx.JSON(http.StatusCreated, res)
}
