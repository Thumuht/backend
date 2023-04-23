/*
Package utils provides utility functions for thumuht.
*/
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

// gql resolver need background information, and they are collected through the http
// request by gin, so we need to save gin context into [Context.context]
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := context.WithValue(ctx.Request.Context(), contextKey("GinContextKey"), ctx)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	}
}

// paired function. eXetract gcontext from context.Context
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

// get me
func GetMe(ctx context.Context) (int, error) {
	gctx, err := GinContextFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return gctx.GetInt("userId"), nil
}
