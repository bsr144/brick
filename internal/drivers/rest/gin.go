package rest

import (
	"brick/internal/config"

	"github.com/gin-gonic/gin"
)

func NewGinServer(appConfig config.App) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	return app
}
