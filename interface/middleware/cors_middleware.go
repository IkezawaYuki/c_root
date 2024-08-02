package middleware

import (
	"github.com/IkezawaYuki/c_root/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:              strings.Split(config.Env.CorsAllowOrigins, ","),
		AllowMethods:              []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:              []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "Access-Control-Allow-Origin", "Referer", "Cookie"},
		AllowCredentials:          true,
		ExposeHeaders:             []string{"Content-Length", "Content-Type", "Set-Cookie"},
		MaxAge:                    -1,
		OptionsResponseStatusCode: 0,
	})
}
