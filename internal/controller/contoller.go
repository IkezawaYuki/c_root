package controller

import (
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/presenter"
	"github.com/IkezawaYuki/c_root/internal/usecase"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type CustomerController struct {
	customerUsecase usecase.CustomerUsecase
	presenter       presenter.Presenter
}

func NewCustomerController(customerUsecase usecase.CustomerUsecase, presenter presenter.Presenter) CustomerController {
	return CustomerController{
		customerUsecase: customerUsecase,
		presenter:       presenter,
	}
}

func (ctr *CustomerController) Login(c echo.Context) error {
	slog.Info("Login is invoked")
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	token, err := ctr.customerUsecase.Login(ctx, &user)
	return c.JSON(ctr.presenter.Generate(err, token))
}

func (ctr *CustomerController) GetCustomer(c echo.Context) error {
	slog.Info("GetCustomer is invoked")
	customerId := c.Param("id")
	ctx := c.Request().Context()
	customer, err := ctr.customerUsecase.GetCustomer(ctx, customerId)
	return c.JSON(ctr.presenter.Generate(err, customer))
}

func (ctr *CustomerController) GetInstagramMedia(c echo.Context) error {
	slog.Info("GetInstagramMedia is invoked")
	customerID := c.Get("customer_id").(string)
	ctx := c.Request().Context()
	mediaList, err := ctr.customerUsecase.GetInstagramMedia(ctx, customerID)
	return c.JSON(ctr.presenter.Generate(err, mediaList))
}

type AdminController struct {
	adminUsecase usecase.AdminUsecase
	presenter    presenter.Presenter
}

func NewAdminController(adminUsecase usecase.AdminUsecase, presenter2 presenter.Presenter) AdminController {
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
