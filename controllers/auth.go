package controllers

import (
	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/lib/controller"
	"github.com/gin-gonic/gin"
)

type Validatable interface {
	ShouldBind(obj interface{}) error
}

type UserService interface {
	CreateUser(ctx Validatable) (interface{}, errors.ApplicationError)
	ChangePassword(user interface{}, ctx Validatable)
}

type MyAuthController struct {
	controller.BaseController
	userService UserService
}

func (ctrl MyAuthController) Signup(c *gin.Context) {
	userDTO, err := ctrl.userService.CreateUser(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
		return
	}
	ctrl.OkResponse(c, controller.AppResponse{Data: userDTO})
}

func (ctrl MyAuthController) ChangePassword(c *gin.Context) {
	user, err := ctrl.GetAuthUser(c)
	if err != nil {
		ctrl.ErrorResponse(c, err)
	}
	ctrl.userService.ChangePassword(user, c)
}
