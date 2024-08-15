package middleware

import (
	"github.com/IkezawaYuki/c_root/internal/presenter"
	"github.com/IkezawaYuki/c_root/internal/service"
	"github.com/labstack/echo/v4"
	"log/slog"
)

func NewBadgeAuthMiddleware(authService service.AuthService, presenter presenter.Presenter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

func NewAdminAuthMiddleware(authService service.AuthService, presenter presenter.Presenter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Info("AdminAuthMiddleware is invoked")
			login, err := authService.IsAdminLogin()
			if !login {
				return c.JSON(presenter.Generate(err, login))
			}
			return next(c)
		}
	}
}

func NewCustomerAuthMiddleware(authService service.AuthService, presenter presenter.Presenter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Info("CustomerAuthMiddleware is invoked")
			token := c.Request().Header.Get("Authorization")
			customerID, err := authService.IsCustomerIsLogin(token)
			if err != nil {
				return c.JSON(presenter.Generate(err, nil))
			}
			slog.Info("login: " + customerID)
			c.Set("customer_id", customerID)
			return next(c)
		}
	}
}
