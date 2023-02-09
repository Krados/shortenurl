package server

import (
	"github.com/Krados/shortenurl/internal/controller"
	"github.com/gin-gonic/gin"
)

func NewHTTPServer(shortenURLController *controller.ShortenURLController) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("api").Group("v1")
	{
		v1.GET("shorten/:code", shortenURLController.Get)
		v1.POST("shorten", shortenURLController.Put)
	}

	return r
}
