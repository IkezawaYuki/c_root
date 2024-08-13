package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	fmt.Println("hello")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/customer/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("id"))
	})
	e.GET("/customer/:id/instagram", func(c echo.Context) error {
		fmt.Println("instagram")
		return c.String(http.StatusOK, c.Param("id"))
	})
	e.Logger.Fatal(e.Start(":1323"))
}
