package handlers

import (
	"fmt"
	"github.com/Anav11/url-shortener/internal/app/storage"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

var URLStorage = storage.GetInstance()

func MainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetHandler(w, r)
	case http.MethodPost:
		PostHandler(w, r)
	default:
		makeResponse(w, "", http.StatusMethodNotAllowed)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if r.URL.Path == "/" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	initialURL := URLStorage.Get(path)

	if initialURL == "" {
		makeResponse(w, "", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", initialURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()
	URLStorage.Add(id, string(body))

	shortURL := fmt.Sprintf("http://%s/%s", r.Host, id)
	makeResponse(w, shortURL, http.StatusCreated)
}

func makeResponse(w http.ResponseWriter, response string, httpStatusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(response))
}
