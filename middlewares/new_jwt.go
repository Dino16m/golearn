package middlewares

import (
	"github.com/dino16m/golearn/adapters"
	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	getAuthUser(claim interface{}) (interface{}, errors.ApplicationError)
}

type JWTAuthMiddleware struct {
	authAdapter adapters.JWTAuthAdapter
	userRepo    UserRepository
}

func (m JWTAuthMiddleware) Authorize(c *gin.Context) {
	authorization, ok := c.Request.Header["Authorization"]
	if !ok {
		errorResponse(c, errors.UnauthorizedError("Unauthorized"))
		return
	}
	token := authorization[len(authorization)-1]
	claims, err := m.authAdapter.GetClaim(token)
	if err != nil {
		errorResponse(c, err)
		return
	}
	uid := claims[types.UserIdClaim]
	user, err := m.userRepo.getAuthUser(uid)
	if err != nil {
		errorResponse(c, errors.UnauthorizedError("User not found"))
	}

	userManager := func() interface{} {
		return user
	}

	c.Set(types.AuthUserContextKey, userManager)
	c.Next()
}

func errorResponse(ctx *gin.Context, err errors.ApplicationError) {
	code, message := err.Resolve()
	ctx.JSON(code, gin.H{
		"status": false,
		"error":  message,
	})
}
