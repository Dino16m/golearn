package container

import (
	"golearn-api-template/config"
	"golearn-api-template/models"

	"gorm.io/gorm"
)

func ProvideDB(cfg config.SuperConfig) *gorm.DB {
	return models.ConnectDataBase(cfg.DatabaseOptions)
}
