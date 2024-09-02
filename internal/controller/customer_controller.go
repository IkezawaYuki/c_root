package controller

import (
	"fmt"
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

// Login godoc
//
//	@Summary		ログイン
//	@Description	顧客のログイン
//	@Accept			application/x-www-form-urlencoded
//	@Security		Token
//	@Param			email			formData	string	false	"Email"
//	@Param			password		formData	string	false	"Password"
//	@Router			/customer/login [post]
func (ctr *CustomerController) Login(c echo.Context) error {
	slog.Info("Login is invoked")
	var user entity.User
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	fmt.Println(user)
	if user.Email == "" || user.Password == "" {
		return c.String(http.StatusBadRequest, "invalid value")
	}
	ctx := c.Request().Context()
	token, err := ctr.customerUsecase.Login(ctx, &user)
	return c.JSON(ctr.presenter.Generate(err, token))
}

// GetCustomer godoc
//
//	@Summary		顧客情報の取得
//	@Description	顧客情報の取得
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Customer ID"
//	@Router			/customer/{id} [get]
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

// FetchAndPost godoc
//
//	@Summary		インスタグラムとWordpressの連携
//	@Description	インスタグラムとWordpressの連携
//	@Produce		json
//	@Security		BearerAuth
//	@Param			customer_id	path	int	true	"Customer ID"
//	@Router			/customer/i/fetch/post [post]
func (ctr *CustomerController) FetchAndPost(c echo.Context) error {
	slog.Info("FetchInstagramMediaFromGraphAPI is invoked")
	customerID := c.Get("customer_id").(int)
	ctx := c.Request().Context()
	err := ctr.customerUsecase.FetchAndPost(ctx, customerID)
	return c.JSON(ctr.presenter.Generate(err, nil))
}
