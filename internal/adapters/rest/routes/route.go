package routes

import (
	"brick/internal/adapters/rest/controllers"
	"brick/internal/adapters/rest/middlewares"

	"github.com/gin-gonic/gin"
)

type RESTRoute struct {
	AuthMiddleware     *middlewares.AuthMiddleware
	UserController     *controllers.UserController
	TransferController *controllers.TransferController
}

func NewRESTRoute(
	userController *controllers.UserController,
	transferController *controllers.TransferController,
	authMiddleware *middlewares.AuthMiddleware,
) *RESTRoute {
	return &RESTRoute{
		AuthMiddleware:     authMiddleware,
		UserController:     userController,
		TransferController: transferController,
	}
}

func (rt *RESTRoute) SetupRoutes(server *gin.Engine) {
	api := server.Group("/api")
	v1 := api.Group("/v1")

	rt.SetupPublicRoutesV1(v1)
	rt.SetupPrivateRoutesV1(v1)
}

func (rt *RESTRoute) SetupPublicRoutesV1(router *gin.RouterGroup) {
	rt.AttachAuthRoutesV1(router)
	rt.AttachTransferPublicRoutesV1(router)
}

func (rt *RESTRoute) SetupPrivateRoutesV1(router *gin.RouterGroup) {
	rt.AttachUserPrivateRoutesV1(router)
	rt.AttachTransferPrivateRoutesV1(router)
}
