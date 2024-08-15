package main

import (
	"fmt"
	"github.com/IkezawaYuki/c_root/config"
	"github.com/IkezawaYuki/c_root/di"
	"github.com/IkezawaYuki/c_root/internal/middleware"
	"github.com/IkezawaYuki/c_root/internal/presenter"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DatabaseUser,
		config.Env.DatabasePass,
		config.Env.DatabaseHost,
		config.Env.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	redisDb := redis.NewClient(&redis.Options{
		Addr: config.Env.RedisAddr,
		DB:   0,
	})

	customerController := di.NewCustomerController(db, redisDb)
	adminController := di.NewAdminController(db, redisDb)
	authService := di.NewAuthService(db, redisDb)
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
