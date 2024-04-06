package routes

import (
	"github.com/gin-gonic/gin"
)

func (rt *RESTRoute) AttachUserPrivateRoutesV1(router *gin.RouterGroup) {
	router.GET("/users/info", rt.AuthMiddleware.VerifyAuth("common"), rt.UserController.GetUserInfoByID)
}
