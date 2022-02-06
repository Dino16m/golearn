package routers

import (
	"github.com/dino16m/golearn/dependencies"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

// Register registers the controller handlers
func Register(c dependencies.ControllersContainer,
	r *gin.Engine, app dependencies.App) {
	apiUserManager := app.ApiUserManager
	webUserManager := app.SessionUserManager
	registerAPIHandlers(c, r, apiUserManager, app)
	registerWebHandlers(c, r, webUserManager, app)
}

func getWrapper(userManager types.AuthUserManager) func(
	gin.HandlerFunc) gin.HandlerFunc {

	return func(h gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set(types.AuthUserContextKey, userManager)
			h(c)
		}
	}
}
