package router

import (
	"github.com/gin-gonic/gin"
)

func ServeThumuht() {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/hello", helloH)
	r.POST("/query", graphqlH())
	r.GET("/", playgroundH())

	r.Run("127.0.0.1:8899")
}
