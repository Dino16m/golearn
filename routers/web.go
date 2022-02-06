package routers

import (
	"github.com/dino16m/golearn/dependencies"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

// Register registers the controller handlers
func registerWebHandlers(c dependencies.ControllersContainer,
	r *gin.Engine, userManager types.AuthUserManager, app dependencies.App) {
	w := getWrapper(userManager)
	web := r.Group("/")
	corsMW := gin.HandlerFunc(app.WebCORS)
	web.Use(corsMW)
	web.Use(app.CSRFMiddleware.Guard())
	web.Use(app.CSRFMiddleware.TokenInjector)

	web.POST("/signup", w(c.AuthController.Signup))
	web.POST("signout", w(c.AuthController.Signout))
	web.POST("signin", w(c.AuthController.Signin))

	web.GET("verify-email", w(c.EmailVerifController.VerifyLink))
	web.GET("verify-password-reset-link", w(c.PasswordResetController.VerifyLink))

	web.POST("send-password-reset-link", w(c.PasswordResetController.SendResetLink))
	web.POST("send-password-reset-code", w(c.PasswordResetController.SendResetCode))
	web.POST("reset-password", w(c.PasswordResetController.ResetPassword))

	web.GET("/", func(c *gin.Context) {
		c.JSON(200, "This is home")
	})
	web.GET("password", func(c *gin.Context) {
		c.HTML(200, "reset-password.html", gin.H{})
	})

	r.LoadHTMLGlob("templates/web/**")

	auth := web.Group("/")
	auth.Use(app.SessionAuthMiddleware.GetHandler())
	{
		auth.POST("authenticate", w(c.AuthController.AuthenticatePassword))
		auth.POST("change-password", w(c.AuthController.ChangePassword))
		auth.GET("users", w(c.AuthController.GetUser))
		auth.POST("resend-email-verification", w(c.EmailVerifController.Resend))
	}
}
