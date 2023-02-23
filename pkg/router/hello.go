package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func helloH(c *gin.Context) {
	s := fmt.Sprintf("Hello World!\nNow Time is %s", time.Now())
	c.String(200, s)
}
