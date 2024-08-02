package router

import (
	"github.com/IkezawaYuki/c_root/interface/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/", "/ping"}}))
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Use(middleware.CorsMiddleware())

	return r
}
