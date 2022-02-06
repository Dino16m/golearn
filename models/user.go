package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Password  string
	Email     string
}

func (u User) GetPassword() string {
	return u.Password
}

func (u User) SetPassword(password string) {
	u.Password = password
}

func (u User) EmailVerified() {

}

func (u User) GetEmail() string {
	return u.Email
}

func (u User) GetId() int {
	return int(u.ID)
}

// GetName returns user name in the format FirstName LastName
func (u User) GetName() (string, string) {
	return u.FirstName, u.LastName
}
