package adapters

import (
	"errors"

	"github.com/dino16m/GinSessionMW/middleware"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

// JwtAuthUserManager used to retrieve a user authenticated using a JWT
type JwtAuthUserManager struct {
	IdentityKey string
	UserRepo    func(identityKey interface{}) types.AuthUser
}

// NewJwtAuthUserManager returns a new instance of JwtAuthUserManager
func NewJwtAuthUserManager(
	identityKey string,
	repo func(identityKey interface{}) types.AuthUser) JwtAuthUserManager {
	return JwtAuthUserManager{identityKey, repo}
}

// GetAuthUser returns the AuthUser authenticated by the JWT middleware
func (manager JwtAuthUserManager) GetAuthUser(
	c *gin.Context,
) (types.AuthUser, error) {
	key, _ := c.Get(manager.IdentityKey)
	if key == nil {
		return nil, errors.New("User not authorized")
	}
	user := manager.UserRepo(key)
	if user == nil {
		return nil, errors.New("User not found")
	}
	return user, nil
}

// SessionAuthUserManager used to retrieve a user authenticated with Session
type SessionAuthUserManager struct {
	middleware *middleware.SessionMiddleware
}

// NewSessionAuthUserManager returns a new instance of a SessionAuthManager
func NewSessionAuthUserManager(
	middleware *middleware.SessionMiddleware,
) SessionAuthUserManager {
	return SessionAuthUserManager{middleware: middleware}
}

// GetAuthUser returns auser authenticated using SessionMiddleware
func (manager SessionAuthUserManager) GetAuthUser(
	c *gin.Context,
) (types.AuthUser, error) {
	user := manager.middleware.GetAuthUser(c)
	if user == nil {
		return nil, errors.New("User not authorized")
	}
	return user.(types.AuthUser), nil
}
