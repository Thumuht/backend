package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func HelloH() gin.HandlerFunc{
	s := fmt.Sprintf("Hello World!\nNow Time is %s", time.Now())
	return func(ctx *gin.Context) {
		ctx.String(200, s)
	} 
}
