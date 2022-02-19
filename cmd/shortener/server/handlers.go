package server

import (
	"fmt"
	"github.com/Anav11/url-shortener/internal/templates"
	"github.com/Anav11/url-shortener/internal/util"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/" {
			fmt.Fprint(w, templates.Form)
		} else {
			initialURL := util.GetURL(r.URL.Path)
			w.Header().Set("Location", initialURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			fmt.Fprint(w, initialURL)
		}
	case http.MethodPost:
		url := r.FormValue("url")
		shortPath := util.URLShortener(url)
		shortURL := fmt.Sprintf("%s/%s", r.Host, shortPath)

		w.WriteHeader(201)
		fmt.Fprint(w, shortURL)
	default:
		http.Error(w, "Bad Request", 400)
	}
}