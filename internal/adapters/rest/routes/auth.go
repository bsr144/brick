package routes

import (
	"github.com/gin-gonic/gin"
)

func (rt *RESTRoute) AttachAuthRoutesV1(router *gin.RouterGroup) {
	router.POST("/auth/login", rt.UserController.LoginUser)
	router.POST("/auth/register", rt.UserController.CreateUser)
	router.POST("/auth/generate-token", rt.UserController.GenerateToken)
}
