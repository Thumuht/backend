package router

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"fmt"
	"time"
)

func GraphqlH(h *handler.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func PlaygroundH() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func HelloH() gin.HandlerFunc {
	s := fmt.Sprintf("Hello World!\nNow Time is %s", time.Now())
	return func(ctx *gin.Context) {
		ctx.String(200, s)
	}
}