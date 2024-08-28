package main

import (
	"github.com/IkezawaYuki/popple/di"
	"github.com/IkezawaYuki/popple/docs"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/middleware"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

//	@title						Recommend Swaggo API
//	@version					1.0
//	@description				This is a recommend_swaggo server
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:1323
//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func main() {
	db := infrastructure.GetMysqlConnection()
	redisCli := infrastructure.GetRedisConnection()
	customerController := di.NewCustomerController(db, redisCli)
	adminController := di.NewAdminController(db, redisCli)
	authService := di.NewAuthService(db, redisCli)
	batchController := di.NewBatchController(db, redisCli)
	pres := presenter.NewPresenter()

	customerAuthMiddleware := middleware.NewCustomerAuthMiddleware(authService, pres)
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(authService, pres)
	badgeAuthMiddleware := middleware.NewBatchAuthMiddleware(authService, pres)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.Group("/api/v1")
	{
		e.POST("/customer/login", customerController.Login)
		e.POST("/admin/login", adminController.Login)

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
	}

	docs.SwaggerInfo.Title = "Popple API"
	docs.SwaggerInfo.Description = "Popple is very very exciting api!!!"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "popple.com"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
