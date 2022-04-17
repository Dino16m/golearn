package controllers

import (
	"net/http"

	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/lib/controller"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	Authenticate(c *gin.Context) (userId interface{}, err errors.ApplicationError)
}
type JWTAuthService interface {
	GetTokenPair(claim map[string]interface{}) (refreshToken string, authToken string)
	GetToken(claim map[string]interface{}) string
	GetClaim(tokenStr string) (map[string]interface{}, errors.ApplicationError)
}

type RefreshTokenPayload struct {
	token string
}

type JWTAuthController struct {
	controller.BaseController
	authenticator Authenticator
	authService   JWTAuthService
}

func (ctrl JWTAuthController) RefreshToken(c *gin.Context) {
	var refresh RefreshTokenPayload
	if err := c.ShouldBind(&refresh); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	claim, err := ctrl.authService.GetClaim(refresh.token)
	if err != nil {
		ctrl.ErrorResponse(c, err)
	}
	freshClaim := map[string]interface{}{
		types.UserIdClaim: claim[types.UserIdClaim],
	}
	token := ctrl.authService.GetToken(freshClaim)
	response := map[string]string{
		"token": token,
	}
	ctrl.OkResponse(c, controller.AppResponse{Data: response})
}

func (ctrl JWTAuthController) GetTokenPair(c *gin.Context) {
	userId, err := ctrl.authenticator.Authenticate(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}
	claim := map[string]interface{}{
		types.UserIdClaim: userId,
	}
	refreshToken, authToken := ctrl.authService.GetTokenPair(claim)
	response := map[string]string{
		"refreshToken": refreshToken,
		"authToken":    authToken,
	}
	ctrl.OkResponse(c, controller.AppResponse{Data: response})
}
