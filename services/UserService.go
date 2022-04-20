package services

import (
	"golearn-api-template/forms"
	"golearn-api-template/models"

	"github.com/dino16m/golearn-core/controller"
	"github.com/dino16m/golearn-core/errors"

	appErrors "golearn-api-template/errors"
)

type UserRepository interface {
	Save(user *models.User)
	Create(user *models.User)
	FindAuthUser(username interface{}) (interface{}, errors.ApplicationError)
}
type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return UserService{userRepo: userRepo}
}
func (service UserService) CreateUser(ctx controller.Validatable) (interface{}, errors.ApplicationError) {
	var userDTO forms.SignUpForm
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		return nil, appErrors.ValidationError(err.Error())
	}
	user := models.NewUser(userDTO.FirstName, userDTO.LastName, userDTO.Email, userDTO.Password)
	service.userRepo.Create(user)
	return user, nil
}

func (service UserService) ChangePassword(authUser interface{}, dto controller.PasswordChangeForm) errors.ApplicationError {
	user := authUser.(models.User)
	ok := user.ComparePassword(dto.OldPassword)
	if !ok {
		return appErrors.ValidationError("Invalid credentials")
	}
	user.SetPassword(dto.NewPassword)
	service.userRepo.Save(&user)
	return nil
}
