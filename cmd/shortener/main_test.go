package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/assert"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/handlers"
	"github.com/Anav11/url-shortener/internal/app/router"
	"github.com/Anav11/url-shortener/internal/app/storage"
)

func TestGetHandler(t *testing.T) {
	c := app.Config{}
	if err := env.Parse(&c); err != nil {
		return
	}

	s := storage.ConstructStorage()
	s.Add("test-id", "https://ya.ru")

	type want struct {
		code int
		contentType string
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "200 http code",
			request: "/test-id",
			want:    want{
				http.StatusTemporaryRedirect,
				"text/plain",
			},
		},
		{
			name:    "404 http code",
			request: "/test-id-fail",
			want:    want{
				http.StatusNotFound,
				"text/plain; charset=utf-8",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.Router(c, s)
			w := httptest.NewRecorder()
			req, err  := http.NewRequest(http.MethodGet, testCase.request, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.Equal(t, testCase.want.code, w.Code)
			assert.NoError(t, err)
		})
	}
}

func TestPostHandler(t *testing.T) {
	c := app.Config{}
	if err := env.Parse(&c); err != nil {
		return
	}

	s := storage.ConstructStorage()

	type want struct {
		code int
		contentType string
	}

	tests := []struct {
		name    string
		want    want
	}{
		{
			name:    "URL added",
			want:    want{
				http.StatusCreated,
				"text/plain",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.Router(c, s)
			w := httptest.NewRecorder()
			req, err  := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}

func TestPostJSONHandler(t *testing.T) {
	c := app.Config{}
	if err := env.Parse(&c); err != nil {
		return
	}

	s := storage.ConstructStorage()

	type want struct {
		code int
		contentType string
	}

	tests := []struct {
		name    string
		want    want
	}{
		{
			name:    "URL JSON added",
			want:    want{
				http.StatusCreated,
				"application/json; charset=utf-8",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := router.Router(c, s)
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(handlers.ShortenerRequestJSON{URL: "https://ya.ru"})
			req, err  := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}
