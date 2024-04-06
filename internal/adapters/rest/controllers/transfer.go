package controllers

import (
	"brick/internal/adapters/dto/request"
	"brick/internal/adapters/dto/response"
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/parser"
	"brick/internal/pkg/serror"
	"brick/internal/usecases/transfer"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransferController struct {
	Log             *logrus.Logger
	TransferUsecase transfer.TransferUseCase
}

func NewTransferController(transferUsecase transfer.TransferUseCase, log *logrus.Logger) *TransferController {
	return &TransferController{
		Log:             log,
		TransferUsecase: transferUsecase,
	}
}

func (c *TransferController) ValidateAccount(ctx *gin.Context) {
	var serr *serror.Error
	accountValidationQueryParams := new(request.ValidateAccount)

	serr = parser.BindQueryParams(ctx, &accountValidationQueryParams)

	if serr != nil {
		c.Log.Errorf("[Transfer][Controller] while ctx.ShouldBindQuery: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	responseValidateAccount, serr := c.TransferUsecase.ValidateAccount(context, accountValidationQueryParams)
	if serr != nil {
		c.Log.Errorf("[Transfer][Controller] while TransferUsecase.ValidateAccount: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	ctx.JSON(http.StatusCreated, response.NewHTTPResponseSuccess(http.StatusCreated, responseValidateAccount))
}

func (c *TransferController) Disburse(ctx *gin.Context) {
	var serr *serror.Error
	disburseRequest := new(request.Disburse)

	userId, err := strconv.Atoi(ctx.GetString("user_id"))
	if err != nil {
		c.Log.Errorf("[User][Controller] while strconv.Atoi: %s", err.Error())
		serror.AbortWithSerror(ctx, http.StatusInternalServerError, 1, constvar.SERVER_INFO_ERROR, err.Error())
		return
	}
	disburseRequest.UserID = userId

	err = ctx.ShouldBind(&disburseRequest)
	if err != nil {
		c.Log.Errorf("[User][Controller] while ctx.ShouldBind: %s", err.Error())
		serror.AbortWithSerror(ctx, http.StatusBadRequest, 1, constvar.BAD_REQUEST_ERROR, err.Error())
		return
	}

	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	responseCreateTransfer, serr := c.TransferUsecase.CreateTransfer(context, disburseRequest)
	if serr != nil {
		c.Log.Errorf("[Transfer][Controller] while TransferUsecase.CreateTransfer: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	transferCallbackRequest := new(request.TransferCallback)
	transferCallbackRequest.TransferID = responseCreateTransfer.TransferID
	// Change the value between "completed" / "pending" / "failed"
	transferCallbackRequest.Status = "Completed"
	go c.TransferUsecase.MockDisburse(context, disburseRequest, transferCallbackRequest)

	ctx.JSON(http.StatusCreated, response.NewHTTPResponseSuccess(http.StatusCreated, responseCreateTransfer))
}

func (c *TransferController) TransferCallback(ctx *gin.Context) {
	var serr *serror.Error
	transferCallbackRequest := new(request.TransferCallback)

	err := ctx.ShouldBind(&transferCallbackRequest)
	if err != nil {
		c.Log.Errorf("[User][Controller] while ctx.ShouldBind: %s", err.Error())
		serror.AbortWithSerror(ctx, http.StatusBadRequest, 1, constvar.BAD_REQUEST_ERROR, err.Error())
		return
	}

	context, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	serr = c.TransferUsecase.UpdateTransferStatus(context, transferCallbackRequest)
	if serr != nil {
		c.Log.Errorf("[Transfer][Controller] while TransferUsecase.CreateTransfer: %s", serr.Error())
		serr.SendAndAbort(ctx)
		return
	}

	ctx.JSON(http.StatusOK, response.NewHTTPResponseSuccess(http.StatusCreated, nil))
}
