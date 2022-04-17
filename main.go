package main

import (
	"context"

	"github.com/dino16m/golearn/config"
	"github.com/dino16m/golearn/dependencies"
	"github.com/dino16m/golearn/models"
	"github.com/dino16m/golearn/registry.go"
	"github.com/dino16m/golearn/routers"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Initialize(".", "yaml")
	db := models.ConnectDataBase()
	userRepo := repositories.NewUserRepository(db)
	repos := dependencies.InitRepos(userRepo)
	mails, err := dependencies.InitMails()
	checkErr(err)
	services, err := dependencies.InitServices(repos, mails)
	checkErr(err)
	App, err := dependencies.InitApp(services, repos)
	checkErr(err)

	r := gin.Default()

	r.Use(App.SessionMw)
	controllers := dependencies.InitControllers(repos, services, App)
	routers.Register(controllers, r, App)
	eventHandlers := dependencies.InitEventHandlers(services, mails)
	registry.RegisterEventHandlers(eventHandlers, App)
	App.Worker.Start(context.Background())

	r.Run()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Setup(
	cfgPath, cfgType string,
	userRepo types.UserRepository, r *gin.Engine) dependencies.App {

	config.Initialize(cfgPath, cfgType)
	repos := dependencies.InitRepos(userRepo)
	mails, err := dependencies.InitMails()
	checkErr(err)
	services, err := dependencies.InitServices(repos, mails)
	checkErr(err)
	App, err := dependencies.InitApp(services, repos)
	checkErr(err)

	r.Use(App.SessionMw)
	controllers := dependencies.InitControllers(repos, services, App)
	routers.Register(controllers, r, App)
	eventHandlers := dependencies.InitEventHandlers(services, mails)
	registry.RegisterEventHandlers(eventHandlers, App)

	return App
}
