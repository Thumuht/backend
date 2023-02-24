package router

import (
	database "backend/pkg/db"
	"github.com/gin-gonic/gin"
)

func ServeThumuht() {
	r := gin.New()
	db, err := database.InitSQLiteDB()
	if err != nil {
		panic("no db")
	}
	database.InitModels(db)

	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/hello", helloH)
	r.POST("/query", graphqlH(db))
	r.GET("/", playgroundH())

	r.Run("127.0.0.1:8899")
}
