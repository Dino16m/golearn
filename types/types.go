package types

import (
	"github.com/dino16m/golearn/models"
	"github.com/gin-gonic/gin"
)

// UserRepository ....
type UserRepository interface {
	GetUserByAuthUsername(interface{}) (AuthUser, error)
	GetUserByAuthId(key interface{}) AuthUser
	CreateUserFromForm(form interface{}) (AuthUser, error)
	Save(interface{})
	GetAllUsers() []models.User // Remove in production
}

// SessionManager ...
type SessionManager interface {
	Login(c *gin.Context, user interface{})
	Logout(c *gin.Context)
}

// Authenticator ...
type Authenticator interface {
	Authenticate(c *gin.Context) (AuthUser, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// AuthUserManager ...
type AuthUserManager interface {
	GetAuthUser(c *gin.Context) (AuthUser, error)
}

// AuthUserContextKey is used by the base controller  and route wrapper
// to fetch or set the AuthUserManager from the request context
const AuthUserContextKey = "AuthUserManager"

// AuthUser ...
type AuthUser interface {
	GetPassword() string
	GetId() int
	EmailVerified()
	SetPassword(string)
	// GetName returns user name in the format FirstName LastName
	GetName() (string, string)
	GetEmail() string
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}
