package middlewares

import (
	"net/http"

	"github.com/dino16m/GinSessionMW/middleware"
	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSessionMw(
	options sessions.Options, userRepo func(interface{}) types.AuthUser,
	sessionFunc func(c *gin.Context) sessions.Session,
) *middleware.SessionMiddleware {

	mw := middleware.New(
		payloadFunc,
		unauthorizedFunc,
		sessions.Options{},
		generifyAuthUserRepo(userRepo),
		func(c *gin.Context) middleware.Session { return sessionFunc(c) },
	)
	return mw
}

func generifyAuthUserRepo(repo func(interface{}) types.AuthUser) middleware.UserRepo {
	return func(i interface{}) interface{} {
		return repo(i)
	}
}
func payloadFunc(user interface{}) interface{} {
	id := user.(types.AuthUser).GetId()
	return id
}

func unauthorizedFunc(c *gin.Context) {
	errors.Unauthorized(
		c, http.StatusUnauthorized, "You are not logged in",
	)
}
