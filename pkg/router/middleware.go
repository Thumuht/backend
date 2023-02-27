package router

import (
	"context"

	"github.com/gin-gonic/gin"
)

// need header information in gql resolver
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := context.WithValue(ctx.Request.Context(), "GinContextKey", ctx)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}