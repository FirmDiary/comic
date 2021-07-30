package main

import (
	"comic/api/controllers"
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	common := mvc.New(app.Party("/common"))
	common.Register(ctx)
	common.Handle(new(controllers.CommonController))

	upload := mvc.New(app.Party("/upload"))
	upload.Register(ctx)
	upload.Handle(new(controllers.UploadController))

	user := mvc.New(app.Party("/user"))
	user.Register(ctx)
	user.Handle(new(controllers.UserController))

	share := mvc.New(app.Party("/invite"))
	share.Register(ctx)
	share.Handle(new(controllers.InviteController))

	app.Run(
		iris.Addr("localhost:8081"),
		//忽略服务器错误
		iris.WithoutServerError(iris.ErrServerClosed),
		//让程序自身尽可能的优化
		iris.WithOptimizations,
	)
}
