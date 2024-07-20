package user

import (
	"github.com/gin-gonic/gin"
)

func (h Handler) response(ctx *gin.Context, statusCode int, msg string) {
	ctx.JSON(statusCode, gin.H{
		"message": msg,
	})
}

func (h Handler) responseWithData(ctx *gin.Context, statusCode int, msg string, data interface{}) {
	ctx.JSON(statusCode, gin.H{
		"message": msg,
		"data":    data,
	})
}

func (h Handler) responseWithError(ctx *gin.Context, statusCode int, err error) {
	ctx.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
