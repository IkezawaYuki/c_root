package router

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/IkezawaYuki/c_root/internal/croot/interface/middleware"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	r := gin.Default()

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/", "/ping"}}))
	r.Use(gin.Recovery())
	r.Use(middleware.CorsMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Env.CRootHost, config.Env.CRootPort),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server Shutdown ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", err.Error())
	}
	slog.Info("Server exiting")
}
