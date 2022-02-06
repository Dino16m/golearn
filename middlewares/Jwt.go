package middlewares

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/dino16m/golearn/config"
	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

var IdentityKey string

type AuthenticatorFunc func(c *gin.Context) (types.AuthUser, error)

func GetJwtMiddleware(
	options config.JwtOptions, authenticator AuthenticatorFunc) (*jwt.GinJWTMiddleware, error) {
	IdentityKey = options.IdentityKey
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           options.Realm,
		Key:             []byte(options.Key),
		Timeout:         options.Timeout,
		MaxRefresh:      options.MaxRefresh,
		IdentityKey:     options.IdentityKey,
		Authenticator:   loginHandler(authenticator),
		PayloadFunc:     payloadHandler,
		Unauthorized:    errors.Unauthorized,
		LoginResponse:   loginResponse,
		RefreshResponse: refreshResponse,
	})
}

func loginHandler(authenticate AuthenticatorFunc) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		user, err := authenticate(c)
		if user != nil {
			return user, nil
		}
		if err == errors.ErrAuthenticationFailed {
			return nil, jwt.ErrFailedAuthentication
		} else {
			return nil, jwt.ErrMissingLoginValues
		}
	}
}

func payloadHandler(data interface{}) jwt.MapClaims {
	if user, ok := data.(types.AuthUser); ok {
		return jwt.MapClaims{
			IdentityKey: user.GetId(),
		}
	}
	return jwt.MapClaims{}
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		},
	})
}

func refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		},
	})
}
