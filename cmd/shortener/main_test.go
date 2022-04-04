package main

import (
	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/router"
	"github.com/Anav11/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
	c := app.Config{
		Host: "http://localhost",
		Port: 8080,
	}

	s := storage.GetInstance()
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
	c := app.Config{
		Host: "http://localhost",
		Port: 8080,
	}
	s := storage.GetInstance()

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
