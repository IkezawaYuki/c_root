package middleware

import (
	"github.com/IkezawaYuki/c_root/di"
	"github.com/IkezawaYuki/c_root/domain/crooterrors"
	"github.com/IkezawaYuki/c_root/interface/controller/presenter"
	"github.com/gin-gonic/gin"
	"github.com/rbcervilla/redisstore/v8"
)

type authMiddleware struct{}

var AuthMiddleware *authMiddleware

func init() {
	AuthMiddleware = &authMiddleware{}
}

func (m authMiddleware) User(ctx *gin.Context) {
	sessionDriver := di.NewSessionDriver()
	client := sessionDriver.GetClient()
	store, err := redisstore.NewRedisStore(ctx, client)
	if err != nil || store == nil {
		ctx.JSON(presenter.Generate(crooterrors.New(crooterrors.UnauthorizedError, err), nil))
		ctx.Abort()
		return
	}

}
