package routers

import (
	"golearn-api-template/container"

	"github.com/gin-gonic/gin"
)

func registerAPI(r *gin.Engine, app container.App) {
	apiGroup := r.Group("/api/v1")
	auth := apiGroup.Group("/auth")
	{
		authCtrl := app.Controllers.AuthController
		authCtrl.RegisterRoutes(auth)

		jwtCtrl := app.Controllers.JWTAuthController
		jwtCtrl.RegisterRoutes(auth)
	}

	secure := apiGroup.Group("/")
	{
		authMW := app.Middlewares.JWTAuthMiddleware
		secure.Use(authMW.Authorize)

	}
}
