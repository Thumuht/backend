/*
Package forum provides thumuht app instance, manages all resources and functionalities.
*/
package forum

import (
	database "backend/pkg/db"
	"backend/pkg/gql/graph"
	"backend/pkg/router"
	"backend/pkg/utils"
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
)

type App struct {
	*gin.Engine                   // router
	DB          *bun.DB           // db instance
	Cache       database.AppCache // cache
	GQL         *handler.Server   // gql server
}

// TODO(wj, low): differ its behavior in accordance with Config..
func NewForum() *App {
	var app App
	var err error

	app.DB, err = database.InitSQLiteDB()
	if err != nil {
		panic("no db")
	}
	database.InitModels(app.DB)
	app.Cache = database.NewAppCache()
	app.Cache.PostLike.SetFlushTarget("post", "post_id", "like", app.DB)
	app.Cache.PostView.SetFlushTarget("post", "post_id", "view", app.DB)
	app.Cache.CommentLike.SetFlushTarget("comment", "comment_id", "like", app.DB)

	app.GQL = handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{
			Resolvers: &graph.Resolver{
				DB:    app.DB,
				Cache: app.Cache,
			},
			Directives: graph.DirectiveRoot{
				Login: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
					gctx, err := utils.GinContextFromContext(ctx)
					if err != nil {
						return nil, fmt.Errorf("cannot get gin context, access denied: %w", err)
					}

					token := gctx.GetHeader("Token")
					if username, ok := app.Cache.Sessions.Get(token); ok {
						gctx.AddParam("appuser", *username)
						return next(ctx)
					}

					return nil, fmt.Errorf("no token %s access denied", token)
				},
			},
		}))

	app.Engine = gin.New()

	SetRouter(&app)

	return &app
}

func SetRouter(app *App) {
	app.Use(gin.Logger(), gin.Recovery(), utils.GinContextToContextMiddleware())

	app.GET("/hello", router.HelloH())

	app.POST("/query", router.GraphqlH(app.GQL))
	app.GET("/", router.PlaygroundH())

	app.StaticFS("/fs", gin.Dir(viper.GetString("fs_route"), true))
}

// Run Forum. BLOCK!!
func (app *App) RunForum(addr string) {
	app.Run(addr)
}
