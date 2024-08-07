package controller

import (
	"github.com/IkezawaYuki/c_root/domain/entity"
	"github.com/IkezawaYuki/c_root/internal/croot/interface/session"
	"github.com/gin-gonic/gin"
)

type customerController struct{}

var CustomerController *customerController

func init() {
	CustomerController = &customerController{}
}

func (c customerController) GetCustomer(ctx *gin.Context) {
	userSession := ctx.MustGet(session.UserSession).(*entity.UserSession)

	di.NewCustomer
}
