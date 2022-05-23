package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/caarlos0/env/v6"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Anav11/url-shortener/internal/app"
	"github.com/Anav11/url-shortener/internal/app/mocks"
	"github.com/Anav11/url-shortener/internal/app/storage"
)

func TestGetHandler(t *testing.T) {
	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		return
	}

	type want struct {
		code int
		contentType string
	}

	tests := []struct {
		mockBehavior func(*mocks.MockRepository, string)
		name 		 string
		request		 string
		want		 want
	}{
		{
			name:    "200 http code",
			request: "/test-id",
			want:    want{
				http.StatusTemporaryRedirect,
				"text/plain",
			},
			mockBehavior: func(s *mocks.MockRepository, id string) {
				s.EXPECT().GetURL(id).Return("/short-test-id", nil)
			},
		},
		{
			name:    "404 http code",
			request: "/test-id-fail",
			want:    want{
				http.StatusNotFound,
				"text/plain; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository, id string) {
				s.EXPECT().GetURL(id).Return("", errors.New(""))
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockRepository(ctrl)
			testCase.mockBehavior(storage, testCase.request)

			handler := Handler{Config: cfg, Storage: storage}

			r := gin.Default()
			r.GET("/:ID", handler.GetHandler)

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
	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		return
	}

	type want struct {
		code int
		contentType string
	}

	tests := []struct {
		name    string
		want    want
		mockBehavior func(*mocks.MockRepository)
	}{
		{
			name: "201 URL created",
			want: want{
				http.StatusCreated,
				"text/plain",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(nil)
			},
		},
		{
			name:    "500 URL creation error",
			want:    want{
				http.StatusInternalServerError,
				"text/plain; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(errors.New(""))
			},
		},
		{
			name:    "409 URL creation duplicate error",
			want:    want{
				http.StatusConflict,
				"text/plain; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(&storage.URLDuplicateError{URL: ""})
				s.EXPECT().GetShortByOriginal(gomock.Any()).Return("test", nil)
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockRepository(ctrl)
			testCase.mockBehavior(storage)

			handler := Handler{Config: cfg, Storage: storage}

			r := gin.Default()
			r.POST("/", handler.PostHandler)

			w := httptest.NewRecorder()
			req, err  := http.NewRequest(http.MethodPost, "/", strings.NewReader("https://test1.ru"))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}

func TestPostJSONHandler(t *testing.T) {
	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		return
	}

	type want struct {
		code 		 int
		contentType  string
	}

	tests := []struct {
		name    string
		want    want
		request		 string
		mockBehavior func(*mocks.MockRepository)
	}{
		{
			name: "201 URL created",
			want: want{
				http.StatusCreated,
				"application/json; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(nil)
			},
		},
		{
			name:    "500 URL creation error",
			want:    want{
				http.StatusInternalServerError,
				"text/plain; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(errors.New(""))
			},
		},
		{
			name:    "409 URL creation duplicate error",
			want:    want{
				http.StatusConflict,
				"application/json; charset=utf-8",
			},
			mockBehavior: func(s *mocks.MockRepository) {
				s.EXPECT().AddURL(gomock.Any()).Return(&storage.URLDuplicateError{URL: ""})
				s.EXPECT().GetShortByOriginal(gomock.Any()).Return("test", nil)
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockRepository(ctrl)
			testCase.mockBehavior(storage)

			handler := Handler{Config: cfg, Storage: storage}

			r := gin.Default()
			r.POST("/api/shorten", handler.PostHandlerJSON)

			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(ShortenerRequestJSON{URL: "https://test3.ru"})
			req, err  := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.want.code, w.Code)
			assert.Equal(t, testCase.want.contentType, w.Header().Get("Content-Type"))
			assert.NoError(t, err)
		})
	}
}
