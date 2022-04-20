package main

import (
	"context"

	"golearn-api-template/container"
	"golearn-api-template/routers"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	app := container.ProvideApp()
	routers.Register(r, app)
	app.Worker.Start(context.Background())

	r.Run()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
