package handlers

import (
	"github.com/Anav11/url-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetHandler(t *testing.T) {
	URLStorage := storage.GetInstance()
	URLStorage.Add("test-id", "https://ya.ru")

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
				"text/plain",
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, testCase.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GetHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, testCase.want.code, result.StatusCode)
		})
	}
}

func TestPostHandler(t *testing.T) {
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
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://ya.ru"))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(PostHandler)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, testCase.want.code, result.StatusCode)
			assert.Equal(t, testCase.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
