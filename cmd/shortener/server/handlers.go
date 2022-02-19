package server

import (
	"fmt"
	"github.com/Anav11/url-shortener/internal/util"
	"io"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			initialURL := util.GetURL(r.URL.Path)
			w.Header().Set("Location", initialURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusInternalServerError)
			return
		}
		shortPath := util.URLShortener(string(body))
		shortURL := fmt.Sprintf("http://%s/%s", r.Host, shortPath)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(201)
		w.Write([]byte(shortURL))
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}