/*
Package router provides http handlers for thumuht app instance.

thumuht app registers these handlers.

*/
package router

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"fmt"
	"time"
)

// graphql server
func GraphqlH(h *handler.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// graphql interactive playground
func PlaygroundH() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(ctx *gin.Context) {
		h.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// hello world page.
//
// check server status.
func HelloH() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(200, fmt.Sprintf("Hello World!\nNow Time is %s", time.Now()))
	}
}
