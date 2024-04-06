package routes

import (
	"github.com/gin-gonic/gin"
)

func (rt *RESTRoute) AttachTransferPrivateRoutesV1(router *gin.RouterGroup) {
	router.GET("/transfer/validate-account", rt.AuthMiddleware.VerifyAuth("api"), rt.TransferController.ValidateAccount)
	router.POST("/transfer", rt.AuthMiddleware.VerifyAuth("api"), rt.TransferController.Disburse)
}

func (rt *RESTRoute) AttachTransferPublicRoutesV1(router *gin.RouterGroup) {
	router.POST("/transfer/callback", rt.TransferController.TransferCallback)
}
