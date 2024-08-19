package controller

import (
	"github.com/IkezawaYuki/popple/internal/domain"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AdminController struct {
	adminUsecase *usecase.AdminUsecase
	presenter    *presenter.Presenter
}

func NewAdminController(adminUsecase *usecase.AdminUsecase, presenter2 *presenter.Presenter) AdminController {
	return AdminController{
		adminUsecase: adminUsecase,
		presenter:    presenter2,
	}
}

func (a *AdminController) RegisterCustomer(c echo.Context) error {
	var customer domain.Customer
	if err := c.Bind(&customer); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := a.adminUsecase.RegisterCustomer(c.Request().Context(), &customer)
	return c.JSON(a.presenter.Generate(err, customer))
}
