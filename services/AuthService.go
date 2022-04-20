package services

import (
	"golearn-api-template/forms"
	"golearn-api-template/models"

	appErrors "golearn-api-template/errors"

	"github.com/dino16m/golearn-core/controller"
	"github.com/dino16m/golearn-core/errors"
)

type AuthService struct {
	userRepo UserRepository
}

type Authenticator interface {
	Authenticate(c controller.Validatable) (userId interface{}, err errors.ApplicationError)
}

func NewAuthService(repo UserRepository) AuthService {
	return AuthService{userRepo: repo}
}

func (service AuthService) Authenticate(c controller.Validatable) (userId interface{}, err errors.ApplicationError) {
	var dto forms.LoginForm
	e := c.ShouldBind(&dto)
	if e != nil {
		return nil, appErrors.ValidationError(e.Error())
	}
	authUser, err := service.userRepo.FindAuthUser(dto.Username)
	user := authUser.(models.User)
	if err != nil {
		return nil, err
	}
	ok := user.ComparePassword(dto.Password)
	if !ok {
		return nil, errors.UnauthorizedError("Invalid login credentials")
	}
	userId = user.ID
	return userId, nil
}
