package controller

import "github.com/gin-gonic/gin"

type batchController struct{}

var BadgeController *batchController

func init() {
	BadgeController = &batchController{}
}

func (b batchController) Execute(ctx *gin.Context) {

}
