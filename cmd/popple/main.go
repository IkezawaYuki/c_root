package main

import (
	"github.com/IkezawaYuki/popple/di"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/middleware"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	db := infrastructure.GetMysqlConnection()
	redisCli := infrastructure.GetRedisConnection()
	customerController := di.NewCustomerController(db, redisCli)
	adminController := di.NewAdminController(db, redisCli)
	authService := di.NewAuthService(db, redisCli)
	pres := presenter.NewPresenter()
	batchController := di.NewBatchController(db, redisCli)

	customerAuthMiddleware := middleware.NewCustomerAuthMiddleware(authService, pres)
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(authService, pres)
	badgeAuthMiddleware := middleware.NewBatchAuthMiddleware(authService, pres)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/customer/login", customerController.Login)

	customerHandler := e.Group("/customer")
	customerHandler.Use(customerAuthMiddleware)
	customerHandler.GET("/:id", func(c echo.Context) error {
		return customerController.GetCustomer(c)
	})
	customerHandler.GET("/:id/instagram", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})
	customerHandler.POST("/:id/facebook_token", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})

	adminHandler := e.Group("/admin")
	adminHandler.Use(adminAuthMiddleware)
	adminHandler.GET("/:id", func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, nil)
	})
	adminHandler.POST("/register/customer", adminController.RegisterCustomer)

	batchHandler := e.Group("/batch")
	batchHandler.Use(badgeAuthMiddleware)
	batchHandler.GET("/execute", batchController.Execute)

	e.Logger.Fatal(e.Start(":1323"))
}
