package forum

import (
	database "backend/pkg/db"
	"backend/pkg/gql/graph"
	"backend/pkg/router"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type App struct {
	*gin.Engine
	DB  *bun.DB
	GQL *handler.Server
}

func NewForum() *App {
	var app App
	var err error

	app.DB, err = database.InitSQLiteDB()
	if err != nil {
		panic("no db")
	}
	database.InitModels(app.DB)

	app.GQL = handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{
				DB:       app.DB,
				Sessions: make(database.Session),
			},
		}))

	app.Engine = gin.New()

	SetRouter(&app)

	return &app
}

func SetRouter(app *App) {
	app.Use(gin.Logger(), gin.Recovery(), router.GinContextToContextMiddleware())

	app.GET("/hello", router.HelloH())

	app.POST("/query", router.GraphqlH(app.GQL))
	app.GET("/", router.PlaygroundH())
}

func (app *App) RunForum(addr string) {
	app.Run(addr)
}
