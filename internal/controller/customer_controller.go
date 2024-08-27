package controller

import (
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type CustomerController struct {
	customerUsecase *usecase.CustomerUsecase
	presenter       *presenter.Presenter
}

func NewCustomerController(customerUsecase *usecase.CustomerUsecase, presenter2 *presenter.Presenter) CustomerController {
	return CustomerController{
		customerUsecase: customerUsecase,
		presenter:       presenter2,
	}
}

func (ctr *CustomerController) Login(c echo.Context) error {
	slog.Info("Login is invoked")
	var user entity.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	token, err := ctr.customerUsecase.Login(ctx, &user)
	return c.JSON(ctr.presenter.Generate(err, token))
}

func (ctr *CustomerController) GetCustomer(c echo.Context) error {
	slog.Info("GetCustomer is invoked")
	customerIdParam := c.Param("id")
	customerId, err := strconv.Atoi(customerIdParam)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	customer, err := ctr.customerUsecase.GetCustomer(ctx, customerId)
	return c.JSON(ctr.presenter.Generate(err, customer))
}

//func (ctr *CustomerController) FetchInstagramMediaFromGraphAPI(c echo.Context) error {
//	slog.Info("FetchInstagramMediaFromGraphAPI is invoked")
//	customerID := c.Get("customer_id").(string)
//	ctx := c.Request().Context()
//	err := ctr.customerUsecase.FetchInstagramMediaFromGraphAPI(ctx, customerID)
//	return c.JSON(ctr.presenter.Generate(err, nil))
//}
