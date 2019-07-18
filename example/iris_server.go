package main

import (
	"github.com/kataras/iris"
	"github.com/sirupsen/logrus"
	"strconv"
)

func main() {
	app := iris.Default()

	app.Get("/hello", func(ctx iris.Context) {
		ctx.WriteString("hello iris!")
	})

	v1 := app.Party("/v1")
	{
		v1.Use(func(ctx iris.Context) {
			logrus.Info("自定义中间件")
			ctx.Next()
		})
		v1.Get("/users/{id:uint64 min(10)}", func(ctx iris.Context) {
			id := ctx.Params().GetUint64Default("id", 10)
			ctx.WriteString(strconv.Itoa(int(id)))
		})

		v1.Get("/orders/{action:string prefix(a_)}", func(ctx iris.Context) {
			action := ctx.Params().Get("action")
			ctx.WriteString(action)
		})
	}

	app.Run(iris.Addr(":8080"))
}
