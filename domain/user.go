package domain

import (
	"togo/config"

	"gorm.io/gorm"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"type:varchar(255);unique"`
	Password string `json:"password" gorm:"type:varchar(255)"`
	Role     Role `json:"role" gorm:"type:varchar(255);default:'user'"` // user, admin
}

type UserService interface {
	Login(user User) (token string, rerr ResponseError)
	SignUp(userSignUp User) (rerr ResponseError)
	ParseToken(token string) (userID uint, rerr ResponseError)
	CreateAdmin(cfg *config.Config)
}

type UserRepository interface {
	GetUserByID(uint) (user User, rerr ResponseError)
	GetUserByEmail(string) (res User, rerr ResponseError)
	CreateUser(user User) (res User, rerr ResponseError)
	// CheckIfExistsByEmail(email string) (isExists bool, rerr ResponseError)
	IsAdmin(userID uint) (isAdmin bool, rerr ResponseError)
}
