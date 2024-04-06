package controllers

import (
	"brick/internal/adapters/dto/request"
	"brick/internal/adapters/dto/response"
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/serror"
	"brick/internal/usecases/user"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserUsecase user.UserUseCase
	Log         *logrus.Logger
}

func NewUserController(userUsecase user.UserUseCase, log *logrus.Logger) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
		Log:         log,
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var serr *serror.Error
	createUserRequest := new(request.CreateUser)

	err := ctx.ShouldBind(&createUserRequest)
	if err != nil {
		c.Log.Errorf("[User][Controller] while ctx.ShouldBind: %s", err.Error())
		serror.AbortWithSerror(ctx, http.StatusBadRequest, 1, constvar.BAD_REQUEST_ERROR, err.Error())
		return
	}

	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	responseCreateUser, serr := c.UserUsecase.CreateUser(context, createUserRequest)
	if serr != nil {
		c.Log.Errorf("[User][Controller] while UserUsecase.CreateUser: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, response.NewHTTPResponseSuccess(http.StatusCreated, responseCreateUser))
}

func (c *UserController) LoginUser(ctx *gin.Context) {
	var serr *serror.Error
	loginUserRequest := new(request.LoginUser)

	err := ctx.ShouldBind(&loginUserRequest)
	if err != nil {
		c.Log.Errorf("[User][Controller] while ctx.ShouldBind: %s", err.Error())
		serror.AbortWithSerror(ctx, http.StatusBadRequest, 1, constvar.BAD_REQUEST_ERROR, err.Error())
		return
	}

	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	responseLoginUser, serr := c.UserUsecase.LoginUser(context, loginUserRequest)
	if serr != nil {
		c.Log.Errorf("[User][Controller] while UserUsecase.LoginUser: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	ctx.JSON(http.StatusOK, response.NewHTTPResponseSuccess(http.StatusOK, responseLoginUser))
}

func (c *UserController) GetUserInfoByID(ctx *gin.Context) {
	var serr *serror.Error
	getInfoRequest := new(request.GetUserInfo)

	userId, err := strconv.Atoi(ctx.GetString("user_id"))
	if err != nil {
		c.Log.Errorf("[User][Controller] while strconv.Atoi: %s", err.Error())
		serror.AbortWithSerror(ctx, http.StatusInternalServerError, 1, constvar.SERVER_INFO_ERROR, err.Error())
		return
	}
	getInfoRequest.ID = userId

	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	responseGetInfoUser, serr := c.UserUsecase.GetUserByID(context, getInfoRequest)
	if serr != nil {
		c.Log.Errorf("[User][Controller] while UserUsecase.GetUserByID: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	ctx.JSON(http.StatusOK, response.NewHTTPResponseSuccess(http.StatusOK, responseGetInfoUser))
}

func (c *UserController) GenerateToken(ctx *gin.Context) {
	var serr *serror.Error
	generateTokenRequest := new(request.GenerateToken)

	clientID, clientSecret, hasAuth := ctx.Request.BasicAuth()
	if !hasAuth {
		c.Log.Errorf("[User][Controller] while ctx.Request.BasicAuth: %s", serr.Error())
		serror.AbortWithSerror(ctx, http.StatusUnauthorized, 1, constvar.UNAUTHORIZED_ERROR, "clientId or clientSecret is missing on basic auth")
		return
	}

	generateTokenRequest.ClientID = clientID
	generateTokenRequest.ClientSecret = clientSecret

	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	responseGenerateToken, serr := c.UserUsecase.GenerateToken(context, generateTokenRequest)
	if serr != nil {
		c.Log.Errorf("[User][Controller] while UserUsecase.GenerateToken: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	ctx.JSON(http.StatusOK, response.NewHTTPResponseSuccess(http.StatusOK, responseGenerateToken))
}
