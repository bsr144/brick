package parser

import (
	"brick/internal/adapters/dto/request"
	"brick/internal/pkg/serror"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Parse incoming JSON body request to a struct,
// 'structName' argument used for error handling information.
func BindJSON(ctx *gin.Context, objectDestination interface{}, structName string) *serror.Error {
	if err := ctx.ShouldBind(&objectDestination); err != nil {
		serror := serror.NewError(http.StatusUnprocessableEntity, 1, fmt.Sprintf("error while parsing JSON body to '%s' struct", structName), err.Error())
		return serror
	}
	return nil
}

// Get integer value of incoming param path variable,
// Example: /users/:id - paramName is 'id' and we want to change it into integer.
func GetIntParam(ctx *gin.Context, paramName string) (int, *serror.Error) {
	val, err := strconv.Atoi(ctx.Param(paramName))
	if err != nil {
		serror := serror.NewError(http.StatusUnprocessableEntity, 1, fmt.Sprintf("error while parsing parameter '%s' to integer", paramName), err.Error())
		return 0, serror
	}
	return val, nil
}

// Get integer value from key-value that previously added via gin `ctx.Set`
func GetIntCtx(ctx *gin.Context, key string) (int, *serror.Error) {
	val, err := strconv.Atoi(ctx.GetString(key))
	if err != nil {
		serror := serror.NewError(http.StatusUnprocessableEntity, 1, fmt.Sprintf("error while parsing context value '%s' to integer", key), err.Error())
		return 0, serror
	}
	return val, nil
}

// Parse incoming query params to a query struct,
// Example: /users?search=test&page=1&size=5 - search, page, and size are the query params we want to parse.
func BindQueryParams(ctx *gin.Context, queries interface{}) *serror.Error {
	err := ctx.ShouldBindQuery(queries)
	if err != nil {
		return serror.NewError(http.StatusUnprocessableEntity, 1, "your request can't be processed", err.Error())
	}

	if val, ok := queries.(*request.Query); ok {
		if val.Page == 0 {
			val.Page = 1
		}

		if val.Size == 0 {
			val.Size = 10
		}
	}

	return nil
}
