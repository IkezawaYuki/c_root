package controller

import (
	"github.com/IkezawaYuki/c_root/internal/domain"
	"github.com/IkezawaYuki/c_root/internal/presenter"
	"github.com/IkezawaYuki/c_root/internal/usecase"
	"github.com/labstack/echo/v4"
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
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	token, err := ctr.customerUsecase.Login(ctx, &user)
	return c.JSON(ctr.presenter.Generate(err, token))
}

func (ctr *CustomerController) GetCustomer(c echo.Context) error {
	customerId := c.Param("id")
	ctx := c.Request().Context()
	customer, err := ctr.customerUsecase.GetCustomer(ctx, customerId)
	return c.JSON(ctr.presenter.Generate(err, customer))
}

//type AuthController struct {
//	usecase.
//	presenter       presenter.Presenter
//}
//
//func NewAuthController(customerSrv service.CustomerService, authSrv service.AuthService, presenter2 presenter.Presenter) AuthController {
//	return AuthController{
//		customerService: customerSrv,
//		authService:     authSrv,
//		presenter:       presenter2,
//	}
//}

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
