package routers

import (
	"github.com/dino16m/golearn/dependencies"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

func registerAPIHandlers(c dependencies.ControllersContainer,
	r *gin.Engine, userManager types.AuthUserManager, app dependencies.App) {
	w := getWrapper(userManager)

	api := r.Group("api")
	corsMW := gin.HandlerFunc(app.APICORS)
	api.Use(corsMW)

	api.POST("/signup", w(c.AuthController.Signup))
	api.POST("signin", w(c.AuthController.Signin))
	api.POST("login", w(app.JwtAuthMW.LoginHandler))
	api.GET("refresh_tokens", app.JwtAuthMW.RefreshHandler)

	api.POST("send-password-reset-link", w(c.PasswordResetController.SendResetLink))
	api.POST("send-password-reset-code", w(c.PasswordResetController.SendResetCode))
	api.POST("reset-password", w(c.PasswordResetController.ResetPassword))

	auth := api.Group("/")
	auth.Use(app.JwtAuthMW.MiddlewareFunc())
	{
		auth.POST("resend-email-verification", w(c.EmailVerifController.Resend))
		auth.POST("change-password", w(c.AuthController.ChangePassword))
		auth.POST("authenticate", w(c.AuthController.AuthenticatePassword))
		auth.GET("users", w(c.AuthController.GetUser))
	}
}
