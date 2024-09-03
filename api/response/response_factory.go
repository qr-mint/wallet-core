package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseFactory struct {
}

func NewResponseFactory() *ResponseFactory {
	return &ResponseFactory{}
}

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func (f *ResponseFactory) InternalServerError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusInternalServerError, Response{Error: message})
}

func (f *ResponseFactory) UnprocessableEntity(ctx *gin.Context, errorMessage string) {
	ctx.JSON(http.StatusUnprocessableEntity, Response{Error: errorMessage})
}

func (f *ResponseFactory) BadRequest(ctx *gin.Context, errorMessage string) {
	ctx.JSON(http.StatusBadRequest, Response{Error: errorMessage})
}

func (f *ResponseFactory) Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{Error: "Unauthorized."})
}

func (f *ResponseFactory) NotFound(ctx *gin.Context, errorMessage string) {
	ctx.JSON(http.StatusNotFound, Response{Error: errorMessage})
}

func (f *ResponseFactory) TooManyRequests(ctx *gin.Context) {
	ctx.JSON(http.StatusTooManyRequests, Response{Error: "Too Many Requests."})
}

func (f *ResponseFactory) Ok(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{Data: data})
}

func (f *ResponseFactory) NoContent(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusNoContent, Response{Data: data})
}
