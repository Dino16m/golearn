package services

import (
	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/forms"
	"github.com/dino16m/golearn/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository types.UserRepository
}

func NewAuthService(repo types.UserRepository) AuthService {
	return AuthService{repo}
}

func (service AuthService) Authenticate(c *gin.Context) (types.AuthUser, error) {
	var loginForm forms.LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		return nil, errors.ErrMissingLoginValues
	}

	username, password := loginForm.Username, loginForm.Password
	user, err := service.repository.GetUserByAuthUsername(username)
	if err != nil {
		return nil, errors.ErrAuthenticationFailed
	}

	authenticated := service.CheckPasswordHash(
		password, user.(types.AuthUser).GetPassword(),
	)

	if authenticated == false {
		return nil, errors.ErrAuthenticationFailed
	}

	return user, nil
}

func (service AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (service AuthService) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
