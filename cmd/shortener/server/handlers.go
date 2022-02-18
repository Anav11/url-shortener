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
			initialUrl := util.GetUrl(r.URL.Path, r.Host)
			w.Header().Set("Location", initialUrl)
			w.WriteHeader(307)
			fmt.Fprint(w, initialUrl)
		}
	case http.MethodPost:
		url := r.FormValue("url")

		if len(url) == 0 {
			http.Error(w, "Empty url", 400)
			return
		}

		shortUrl := util.UrlShortener(url, r.Host)
		w.WriteHeader(201)
		fmt.Fprint(w, shortUrl)
	default:
		http.Error(w, "Bad Request", 400)
	}
}