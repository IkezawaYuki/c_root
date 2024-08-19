package controller

import (
	"github.com/IkezawaYuki/popple/di"
	"github.com/IkezawaYuki/popple/internal/croot/domain/entity"
	"github.com/IkezawaYuki/popple/internal/croot/interfaces/presenter"
	"github.com/IkezawaYuki/popple/internal/croot/interfaces/session"
	"github.com/gin-gonic/gin"
)

type customerController struct{}

var CustomerController *customerController

func init() {
	CustomerController = &customerController{}
}

func (c customerController) GetInstagram(ctx *gin.Context) {
	userSession := ctx.MustGet(session.UserSession).(*entity.UserSession)
	s := di.NewCustomerService()
	medias, err := s.GetInstagram(ctx, userSession.UserID)
	ctx.JSON(presenter.Generate(err, medias))
}
