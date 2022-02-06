package controllers

import (
	"net/http"
	"reflect"

	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/events"
	"github.com/dino16m/golearn/forms"
	"github.com/dino16m/golearn/lib/controller"
	"github.com/dino16m/golearn/lib/event"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	sessionManager types.SessionManager
	authenticator  types.Authenticator
	userRepo       types.UserRepository
	form           reflect.Type
	dispatcher     event.Dispatcher
	controller.BaseController
}

func NewAuthController(
	sessionManager types.SessionManager,
	authenticator types.Authenticator,
	userRepo types.UserRepository,
	form interface{},
	dispatcher event.Dispatcher,
) AuthController {
	return AuthController{
		sessionManager: sessionManager,
		authenticator:  authenticator,
		userRepo:       userRepo,
		form:           reflect.TypeOf(form),
		dispatcher:     dispatcher,
	}

}

func (controller AuthController) Signin(c *gin.Context) {
	user, err := controller.authenticator.Authenticate(c)
	if err != nil {
		appErr := err.(errors.AppError)
		c.JSON(appErr.Code, gin.H{
			"status": false,
			"error":  appErr.Message,
		})

	} else {
		controller.sessionManager.Login(c, user)
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"error":  nil,
		})
	}
}
func (controller AuthController) GetUser(c *gin.Context) {
	users := controller.userRepo.GetAllUsers()
	c.JSON(http.StatusOK, users)
}
func (controller AuthController) Signout(c *gin.Context) {
	controller.sessionManager.Logout(c)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"error":  nil,
	})
}

// Signup gin handler for signup
func (controller AuthController) Signup(c *gin.Context) {
	signupFormPointer := reflect.New(controller.form)
	signupFormIface := signupFormPointer.Interface()
	if err := c.ShouldBind(signupFormIface); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	password := reflect.Indirect(signupFormPointer).FieldByName("Password")
	hashedPassword, err := controller.authenticator.HashPassword(
		password.String(),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Invalid password",
		})
		return
	}
	passwordValue := reflect.ValueOf(hashedPassword)
	reflect.Indirect(signupFormPointer).FieldByName("Password").Set(passwordValue)

	user, usererr := controller.userRepo.CreateUserFromForm(signupFormIface)
	if usererr == nil {
		controller.sessionManager.Login(c, user)
		c.JSON(http.StatusCreated, gin.H{
			"status": true,
			"error":  nil,
		})
		firstName, lastName := user.GetName()

		payload := map[string]interface{}{
			"firstName": firstName,
			"lastName":  lastName,
			"email":     user.GetEmail(),
			"id":        user.GetId(),
			"baseURL":   controller.GetBaseURL(c),
		}
		controller.dispatcher.Dispatch(events.UserCreated, payload)
	} else {
		c.JSON(http.StatusNotModified, gin.H{
			"status": false,
			"error":  "User creation failed",
		})
	}
}

func (controller AuthController) AuthenticatePassword(c *gin.Context) {
	var passwordCheckForm forms.PasswordCheckForm
	if err := c.ShouldBind(&passwordCheckForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	user, err := controller.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  "User not found",
		})
		return
	}
	userPassword := user.GetPassword()
	isEqual := controller.authenticator.CheckPasswordHash(
		passwordCheckForm.Password, userPassword,
	)

	status := isEqual
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"error":  nil,
	})
}

func (controller AuthController) ChangePassword(c *gin.Context) {
	var form forms.PasswordChangeForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	user, err := controller.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  "User not found",
		})
		return
	}
	userPassword := user.GetPassword()
	isEqual := controller.authenticator.CheckPasswordHash(
		form.OldPassword, userPassword,
	)
	if isEqual == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
			"error":  "Incorrect password",
		})
		return
	}
	hashedPassword, _ := controller.authenticator.HashPassword(
		form.NewPassword,
	)
	user.SetPassword(hashedPassword)
	controller.userRepo.Save(user)
}
