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
	"github.com/99designs/gqlgen/graphql/handler/transport"
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
					if userId, ok := app.Cache.Sessions.Get(token); ok {
						// let id be a string
						gctx.AddParam("userId", fmt.Sprintf("%d", userId))
						return next(ctx)
					}

					return nil, fmt.Errorf("no token %s access denied", token)
				},
			},
		}))

	app.GQL.AddTransport(&transport.Websocket{})

	app.Engine = gin.New()

	SetRouter(&app)

	// set up a system account
	setUpSysAcc(&app)

	return &app
}

func setUpSysAcc(app *App) {
	// create a system account
	sys := database.User{
		Nickname: "论坛小助手",
	}
	_, err := app.DB.NewInsert().Model(&sys).Exec(context.Background())
	if err != nil {
		panic("cannot create system account")
	}
}

func SetRouter(app *App) {
	app.Use(gin.Logger(), gin.Recovery(), utils.GinContextToContextMiddleware())

	app.GET("/hello", router.HelloH())

	app.POST("/query", router.GraphqlH(app.GQL))
	app.GET("/query", router.GraphqlH(app.GQL))
	app.GET("/", router.PlaygroundH())

	app.StaticFS("/fs", gin.Dir(viper.GetString("fs_route"), true))
	app.POST("/upload", func(ctx *gin.Context) {
		ctx.String(200, "upload")
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.String(200, "no file")
			return
		}
		ctx.String(200, file.Filename)
		err = ctx.SaveUploadedFile(file, viper.GetString("fs_route")+"/"+file.Filename)
		if err != nil {
			ctx.String(200, "save error")
			return
		}
	})
}

// Run Forum. BLOCK!!
func (app *App) RunForum(addr string) {
	app.Run(addr)
}
