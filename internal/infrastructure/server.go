package infrastructure

import (
	"github.com/IkezawaYuki/c_root/di"
	"github.com/IkezawaYuki/c_root/internal/middleware"
	"github.com/IkezawaYuki/c_root/internal/presenter"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Run() {
	db := GetMysqlConnection()
	redisCli := GetRedisConnection()
	customerController := di.NewCustomerController(db, redisCli)
	adminController := di.NewAdminController(db, redisCli)
	authService := di.NewAuthService(db, redisCli)
	pres := presenter.NewPresenter()

	customerAuthMiddleware := middleware.NewCustomerAuthMiddleware(authService, pres)
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(authService, pres)
	badgeAuthMiddleware := middleware.NewBadgeAuthMiddleware(authService, pres)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/customer/login", customerController.Login)

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

	badgeHandler := e.Group("/badge")
	badgeHandler.Use(badgeAuthMiddleware)
	badgeHandler.GET("/:id", func(c echo.Context) error {
		return c.JSON(http.StatusNotImplemented, nil)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
