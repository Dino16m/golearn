package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDataBase() *gorm.DB {
	database, err := gorm.Open(
		mysql.Open("UZU:uzu1@tcp(127.0.0.1:3306)/gonic?parseTime=true"),
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
