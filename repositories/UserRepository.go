package repositories

import (
	"errors"

	appErrors "golearn-api-template/errors"

	coreErrors "github.com/dino16m/golearn-core/errors"

	"golearn-api-template/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (repo UserRepository) Create(user *models.User) {
	repo.db.Create(user)
}

func (repo UserRepository) Save(user *models.User) {
	repo.db.Save(user)
}

func (repo UserRepository) FindAuthUser(username interface{}) (interface{}, coreErrors.ApplicationError) {
	var user models.User
	res := repo.db.Where("email = ?", username).First(&user)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &models.User{}, appErrors.NotFound("User not found")
	}
	return &user, nil
}
