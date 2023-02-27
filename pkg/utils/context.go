package utils

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// need header information in gql resolver
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := context.WithValue(ctx.Request.Context(), contextKey("GinContextKey"), ctx)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(contextKey("GinContextKey"))
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil

}
