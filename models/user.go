package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	password  string
	Email     string
}

func NewUser(firstName, lastName, email, password string) *User {
	return &User{
		FirstName: firstName, LastName: lastName, Email: email, password: hashPassword(password),
	}
}

func (u *User) SetPassword(password string) {
	u.password = hashPassword(password)
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}
