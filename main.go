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

	uploadParty := app.Party("/upload")
	upload := mvc.New(uploadParty)
	upload.Register(ctx)
	upload.Handle(new(controllers.UploadController))

	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx)
	user.Handle(new(controllers.UserController))

	app.Run(
		iris.Addr("localhost:8081"),
		//忽略服务器错误
		iris.WithoutServerError(iris.ErrServerClosed),
		//让程序自身尽可能的优化
		iris.WithOptimizations,
	)
}
