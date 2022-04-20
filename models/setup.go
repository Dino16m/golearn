package models

import (
	"golearn-api-template/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDataBase(cfg config.DatabaseOptions) *gorm.DB {
	database, err := gorm.Open(
		mysql.Open(cfg.URL),
		&gorm.Config{},
	)

	if err != nil {
		panic("Failed to connect to database!")
	}

	migrateModels(database)
	return database
}

func migrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
