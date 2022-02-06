package controller

import (
	"errors"
	"fmt"

	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

type authUserManager interface {
	GetAuthUser(c *gin.Context) (interface{}, error)
}

// A BaseController is the
// base struct implementing behaviour universal to controllers
type BaseController struct {
}

// GetAuthUser returns the authenticated user interface and a nil error
// if such user exists.
// It returns a nil user and an error if the user does not exist or if
// no auth manager was registered.
func (b BaseController) GetAuthUser(c *gin.Context) (types.AuthUser, error) {

	manager, exists := c.Get(types.AuthUserContextKey)
	if exists == false {
		return nil, errors.New("No user manager provided")
	}
	userManager, ok := manager.(types.AuthUserManager)
	if ok == false {
		return nil, errors.New("Invalid user manager")
	}
	return userManager.GetAuthUser(c)
}

// GetBaseURL return the fully qualified url to the root of the app, from the
// request url
func (b BaseController) GetBaseURL(c *gin.Context) string {
	scheme := "http"
	host := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	return baseURL
}
