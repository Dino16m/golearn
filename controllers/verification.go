package controllers

import (
	"errors"
	"net/http"

	"github.com/dino16m/golearn/lib/controller"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
)

// EmailVerificationService service interface expected by the controller
type EmailVerificationService interface {
	VerifyAndGetUserID(requestURL string) (int, error)
	SendVerificationMail(name string, id int, email, link string) error
}

// EmailVerifController ...
type EmailVerifController struct {
	verifService EmailVerificationService
	userRepo     types.UserRepository
	controller.BaseController
}

// NewEmailVerifController construct the controller
func NewEmailVerifController(
	verifService EmailVerificationService,
	repo types.UserRepository) EmailVerifController {
	return EmailVerifController{
		userRepo:     repo,
		verifService: verifService,
	}
}

// VerifyLink ...
func (evc EmailVerifController) VerifyLink(c *gin.Context) {
	url := c.Request.URL.String()
	userID, err := evc.verifService.VerifyAndGetUserID(url)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid link data"))
		return
	}
	user := evc.userRepo.GetUserByAuthId(userID)
	user.EmailVerified()
	evc.userRepo.Save(user)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// Resend ...
func (evc EmailVerifController) Resend(c *gin.Context) {
	user, err := evc.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status": false,
			"error":  "User not authhorized",
		})
		return
	}
	email := user.GetEmail()
	id := user.GetId()
	firstName, _ := user.GetName()

	err = evc.verifService.SendVerificationMail(firstName, id, email, evc.GetBaseURL(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "Mailing failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
