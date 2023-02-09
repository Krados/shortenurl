package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Krados/shortenurl/internal/service/cache"
	"github.com/Krados/shortenurl/internal/service/shortener"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupRouter(controller *ShortenURLController) *gin.Engine {
	r := gin.Default()
	r.GET("/shorten/:code", controller.Get)
	r.POST("/shorten", controller.Put)
	return r
}

func TestGet_CacheInternalError(t *testing.T) {
	// init controller
	code := "abcdef"
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", errors.New("some bad error"))
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// set up http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shorten/"+code, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGet_CacheHit(t *testing.T) {
	// init controller
	code := "abcdef"
	rURL := "https://www.google.com"
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return(rURL, nil)
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shorten/"+code, nil)
	router.ServeHTTP(w, req)

	// test
	b := w.Body.Bytes()
	var resp struct {
		LongURL string `json:"long_url"`
	}
	err := json.Unmarshal(b, &resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, rURL, resp.LongURL)
}

func TestGet_CodeNotFound(t *testing.T) {
	// init controller
	code := "abcdef"
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", cache.ErrKeyNotFound)
	ms.On("Get", mock.Anything).Return(nil, shortener.ErrCodeNotFound)
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shorten/"+code, nil)
	router.ServeHTTP(w, req)

	// test
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGet_ShortenerInternalError(t *testing.T) {
	// init controller
	code := "abcdef"
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", cache.ErrKeyNotFound)
	ms.On("Get", mock.Anything).Return(nil, errors.New("some bad error"))
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shorten/"+code, nil)
	router.ServeHTTP(w, req)

	// test
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGet_Success(t *testing.T) {
	// init controller
	code := "abcdef"
	rURL := "https://www.google.com"
	shortURL := shortener.ShortURL{
		ID:      1,
		Code:    code,
		HashURL: "HashURL",
		LongURL: rURL,
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", cache.ErrKeyNotFound)
	mc.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	ms.On("Get", mock.Anything).Return(&shortURL, nil)
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/shorten/"+code, nil)
	router.ServeHTTP(w, req)

	// test
	b := w.Body.Bytes()
	var resp struct {
		LongURL string `json:"long_url"`
	}
	err := json.Unmarshal(b, &resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, rURL, resp.LongURL)
}

func TestPut_EmptyLongURL(t *testing.T) {
	// init controller
	buf := bytes.NewBufferString(`{"long_url":""}`)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten", buf)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// test
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPut_MissingFromBody(t *testing.T) {
	// init controller
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// test
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPut_CacheInternalError(t *testing.T) {
	// init controller
	buf := bytes.NewBufferString(`{"long_url":"https://www.google.com"}`)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", errors.New("some bad error"))
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten", buf)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// test
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPut_CacheHit(t *testing.T) {
	// init controller
	code := "abcdef"
	buf := bytes.NewBufferString(`{"long_url":"https://www.google.com"}`)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return(code, nil)
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten", buf)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// test
	var resp struct {
		Code string `json:"code"`
	}
	b := w.Body.Bytes()
	err := json.Unmarshal(b, &resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, code, resp.Code)
}

func TestPut_InvalidURL(t *testing.T) {
	// init controller
	buf := bytes.NewBufferString(`{"long_url":"wwwgooglecom"}`)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", cache.ErrKeyNotFound)
	ms.On("Put", mock.Anything).Return(nil, shortener.ErrInValidURL)
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten", buf)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// test
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPut_ShortenerInternalError(t *testing.T) {
	// init controller
	buf := bytes.NewBufferString(`{"long_url":"wwwgooglecom"}`)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", cache.ErrKeyNotFound)
	ms.On("Put", mock.Anything).Return(nil, errors.New("some bad error"))
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten", buf)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// test
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestPut_Success(t *testing.T) {
	// init controller
	code := "abcdef"
	rURL := "https://www.google.com"
	shortURL := shortener.ShortURL{
		ID:      1,
		Code:    code,
		HashURL: "HashURL",
		LongURL: rURL,
	}
	buf := bytes.NewBufferString(`{"long_url":"https://www.google.com"}`)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ms := shortener.MockShortener{}
	mc := cache.MockCache{}
	mc.On("Get", mock.Anything).Return("", cache.ErrKeyNotFound)
	mc.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	ms.On("Put", mock.Anything).Return(&shortURL, nil)
	controller := NewShortenURLController(sugar, &ms, &mc)

	// init router
	router := setupRouter(controller)

	// setup http test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/shorten", buf)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// test
	var resp struct {
		Code string `json:"code"`
	}
	b := w.Body.Bytes()
	err := json.Unmarshal(b, &resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, code, resp.Code)
}
