package controller

import (
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/labstack/echo/v4"
	"log/slog"
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

// RegisterCustomer godoc
//
//	@Summary		顧客情報の作成
//	@Description	顧客情報の作成
//	@Produce		json
//	@Security		Token
//	@Param			name			formData	string	false	"Name"
//	@Param			password		formData	string	false	"Password"
//	@Param			email			formData	string	false	"Email"
//	@Param			wordpress_url	formData	string	false	"WordPress URL"
//	@Router			/admin/register/customer [post]
func (a *AdminController) RegisterCustomer(c echo.Context) error {
	slog.Info("RegisterCustomer is invoked")
	var customer entity.Customer
	customer.Name = c.FormValue("name")
	customer.Password = c.FormValue("password")
	customer.Email = c.FormValue("email")
	customer.WordpressURL = c.FormValue("wordpress_url")
	if customer.Name == "" || customer.Password == "" || customer.Email == "" || customer.WordpressURL == "" {
		return c.JSON(http.StatusBadRequest, "invalid value")
	}
	err := a.adminUsecase.RegisterCustomer(c.Request().Context(), &customer)
	return c.JSON(a.presenter.Generate(err, customer))
}

// Login godoc
//
//	@Summary		ログイン
//	@Description	管理者のログイン
//	@Accept			application/x-www-form-urlencoded
//	@Security		Token
//	@Param			email			formData	string	false	"Email"
//	@Param			password		formData	string	false	"Password"
//	@Router			/admin/login [post]
func (a *AdminController) Login(c echo.Context) error {
	var user entity.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	token, err := a.adminUsecase.Login(c.Request().Context(), &user)
	return c.JSON(a.presenter.Generate(err, token))
}
