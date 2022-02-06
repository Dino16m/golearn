package controllers

import (
	"net/http"

	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/lib/controller"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

type PasswordResetService interface {
	VerifyAndGetUserID(code int, payload string) (int, error)
	SendPasswordResetLink(name string, id int, email, baseURL string) error
	SendPasswordResetCode(name string, id int, email string) (string, error)
	GetURLClaims(requestURL string) (int, string)
}
type requestForm struct {
	Username string `json:"username" form:"username"`
}
type resetForm struct {
	Code     int    `json:"code" form:"code"`
	Payload  string `json:"signature" form:"signature"`
	Password string `json:"password" form:"password"`
}

type mailerParams struct {
	name  string
	email string
	id    int
}
type PasswordResetController struct {
	resetService PasswordResetService
	userRepo     types.UserRepository
	controller.BaseController
	authService types.Authenticator
	appName     string
}

func NewPasswordResetController(resetService PasswordResetService,
	userRepo types.UserRepository, authService types.Authenticator,
	appName string) PasswordResetController {
	return PasswordResetController{
		resetService: resetService,
		userRepo:     userRepo,
		appName:      appName,
		authService:  authService,
	}
}

func (prc PasswordResetController) SendResetLink(c *gin.Context) {
	params, err := prc.getMailerParams(c)
	if err != nil {
		appErr := err.(errors.AppError)
		c.JSON(appErr.Code, gin.H{
			"status": false,
			"error":  appErr.Message,
		})
		return
	}
	err = prc.resetService.SendPasswordResetLink(
		params.name, params.id, params.email, prc.GetBaseURL(c),
	)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
		})
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"status": false,
			"error":  "unable to send mail please contact support",
		})
	}
}

func (prc PasswordResetController) SendResetCode(c *gin.Context) {
	params, err := prc.getMailerParams(c)
	if err != nil {
		appErr := err.(errors.AppError)
		c.JSON(appErr.Code, gin.H{
			"status": false,
			"error":  appErr.Message,
		})
		return
	}
	signature, err := prc.resetService.SendPasswordResetCode(
		params.name, params.id, params.email,
	)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"data":   gin.H{"signature": signature},
		})
	} else {
		c.JSON(http.StatusAccepted, gin.H{
			"status": false,
			"error":  "unable to send mail please contact support",
		})
	}
}
func (prc PasswordResetController) getMailerParams(c *gin.Context) (mailerParams, error) {
	var form requestForm
	if err := c.ShouldBind(&form); err != nil {
		return mailerParams{}, errors.AppError{
			Code: http.StatusBadRequest, Message: "username required",
		}
	}
	username := form.Username
	user, err := prc.userRepo.GetUserByAuthUsername(username)
	if err != nil {
		return mailerParams{}, errors.AppError{
			Code: http.StatusOK, Message: "username required"}
	}
	email := user.GetEmail()
	id := user.GetId()
	firstName, _ := user.GetName()
	return mailerParams{firstName, email, id}, nil
}

func (prc PasswordResetController) VerifyLink(c *gin.Context) {
	reqURL := c.Request.URL.String()
	code, payload := prc.resetService.GetURLClaims(reqURL)
	c.HTML(http.StatusOK, "reset-password.html", gin.H{
		"code":      code,
		"signature": payload,
		"appName":   prc.appName,
	})
}

func (prc PasswordResetController) ResetPassword(c *gin.Context) {
	var form resetForm
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	userID, err := prc.resetService.VerifyAndGetUserID(form.Code, form.Payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}
	user := prc.userRepo.GetUserByAuthId(userID)
	hashedPassword, _ := prc.authService.HashPassword(form.Password)
	user.SetPassword(hashedPassword)
	prc.userRepo.Save(user)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
