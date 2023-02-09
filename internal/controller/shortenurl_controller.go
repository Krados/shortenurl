package controller

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/Krados/shortenurl/internal/service/cache"
	"github.com/Krados/shortenurl/internal/service/shortener"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShortenURLController struct {
	shortener      shortener.Shortener
	sCache         cache.Cache
	cacheGetFormat string
	cachePutFormat string
	logger         *zap.SugaredLogger
}

func NewShortenURLController(logger *zap.SugaredLogger, shortener shortener.Shortener, sCache cache.Cache) *ShortenURLController {
	return &ShortenURLController{
		shortener:      shortener,
		sCache:         sCache,
		cacheGetFormat: "ShortenURLGET_%v",
		cachePutFormat: "ShortenURLPUT_%v",
		logger:         logger,
	}
}

func (s *ShortenURLController) getCacheGetKey(key string) string {
	return fmt.Sprintf(s.cacheGetFormat, key)
}

func (s *ShortenURLController) getCachePutKey(key string) string {
	return fmt.Sprintf(s.cachePutFormat, key)
}

func (s *ShortenURLController) Get(c *gin.Context) {
	code := c.Param("code")
	cacheKey := s.getCacheGetKey(code)
	// check cache exist or not
	cRes, err := s.sCache.Get(cacheKey)
	if err != nil && err != cache.ErrKeyNotFound {
		s.logger.Errorf("s.sCache.Get err:%s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	// if exist then return
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"long_url": cRes,
		})
		return
	}

	// get from the database
	su, err := s.shortener.Get(code)
	if err != nil {
		if err == shortener.ErrCodeNotFound {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		s.logger.Errorf("s.shortener.Get err:%s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	// set data to cache
	err = s.sCache.Set(cacheKey, su.LongURL, time.Minute)
	if err != nil {
		s.logger.Errorf("s.sCache.Set err:%v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"long_url": su.LongURL,
	})
}

func (s *ShortenURLController) Put(c *gin.Context) {
	var req struct {
		LongURL string `json:"long_url"`
	}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if req.LongURL == "" {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	hashURL := fmt.Sprintf("%x", sha256.Sum256([]byte(req.LongURL)))
	cacheKey := s.getCachePutKey(hashURL)
	// check cache exist or not
	cRes, err := s.sCache.Get(cacheKey)
	if err != nil && err != cache.ErrKeyNotFound {
		s.logger.Errorf("s.sCache.Get err:%s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	// if exist then return
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": cRes,
		})
		return
	}

	// set data to database
	su, err := s.shortener.Put(req.LongURL)
	if err != nil {
		if err == shortener.ErrInValidURL {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		s.logger.Errorf("s.shortener.Put err:%s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	// set data to cache
	err = s.sCache.Set(cacheKey, su.Code, time.Minute)
	if err != nil {
		s.logger.Errorf("s.sCache.Set err:%v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": su.Code,
	})
}
