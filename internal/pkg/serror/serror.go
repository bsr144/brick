package serror

import (
	"brick/internal/adapters/dto/response"
	"brick/internal/pkg/constvar"
	"runtime"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code     int           `json:"code"`
	Message  string        `json:"message"`
	Info     string        `json:"info"`
	Detail   string        `json:"detail"`
	Location ErrorLocation `json:"location"`
}

type ErrorLocation struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

func (e *Error) Error() string {
	return e.Detail
}

func (e *Error) SendAndAbort(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(e.Code, response.NewHTTPResponseError(e.Code, e))
}

func AbortWithSerror(ctx *gin.Context, code, skip int, info, detail string) {
	serr := NewError(code, skip, info, detail)
	ctx.AbortWithStatusJSON(code, response.NewHTTPResponseError(code, serr))
}

// 'Info' parameter is used for more friendly error message
// 'Detail' parameter is used for more detailed error message <usually comes from the package or function used directly, e.g err.Error()>
func NewError(code, skip int, info, detail string) *Error {
	_, file, line, _ := runtime.Caller(1 + skip)

	err := &Error{
		Code:    code,
		Message: statusMessage(code),
		Info:    info,
		Detail:  detail,
		Location: ErrorLocation{
			File: file,
			Line: line,
		},
	}

	return err
}

// This function copied from fiber framework utils package
// Thanks to fiber utils package, could be accessed via: https://github.com/gofiber/fiber/blob/v2.51.0/utils
func statusMessage(code int) string {
	if code < constvar.StatusCodeMessageMin || code > constvar.StatusCodeMessageMax {
		return "status code out of scope (less than 400 OR more than 511)"
	}
	return constvar.StatusCodeMessage[code]
}
