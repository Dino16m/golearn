package repositories

import (
	"github.com/dino16m/golearn/forms"
	"github.com/dino16m/golearn/models"
	"github.com/dino16m/golearn/types"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db}
}

func (repo UserRepository) GetUserByAuthId(key interface{}) types.AuthUser {
	var user models.User
	result := repo.db.First(&user, key)
	if result.Error != nil {
		return nil
	}
	return user
}

func (repo UserRepository) GetUserByAuthUsername(key interface{}) (types.AuthUser, error) {
	var user models.User
	email := key.(string)
	result := repo.db.Where(&models.User{Email: email}).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (repo UserRepository) CreateUserFromForm(
	form interface{},
) (types.AuthUser, error) {
	signupForm := form.(*forms.SignUpForm)
	user := models.User{
		Email:     signupForm.Email,
		Password:  signupForm.Password,
		FirstName: signupForm.FirstName,
		LastName:  signupForm.LastName,
	}
	result := repo.db.Create(&user)
	return user, result.Error
}

func (repo UserRepository) GetAllUsers() []models.User {
	var users []models.User
	repo.db.Find(&users)
	return users
}

func (repo UserRepository) Save(user interface{}) {
	repo.db.Save(user)
}
