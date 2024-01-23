package repository

import (
	"errors"
	"fmt"
	"togo/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	conn *gorm.DB
}

func NewUserReposity(conn *gorm.DB) domain.UserRepository {
	return &userRepository{conn: conn}
}

func (uRepo *userRepository) IsAdmin(userID uint) (bool, domain.ResponseError) {
	user, rerr := uRepo.GetUserByID(userID)
	if rerr != nil {
		return false, rerr
	}

	return user.Role == domain.AdminRole, nil
}

func (uRepo *userRepository) GetUserByID(userID uint) (domain.User, domain.ResponseError) {
	var user domain.User

	result := uRepo.conn.Where("id = ?", userID).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.User{}, domain.NewReponseError(
			domain.ErrorNotFound, result.Error.Error(),
		)
	}
	if result.Error != nil {
		return domain.User{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return user, nil
}

func (uRepo *userRepository) GetUserByEmail(email string) (domain.User, domain.ResponseError) {
	var user domain.User

	result := uRepo.conn.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.User{}, domain.NewReponseError(
			domain.ErrorNotFound, result.Error.Error(),
		)
	}
	if result.Error != nil {
		return domain.User{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return user, nil
}

func (uRepo *userRepository) CreateUser(user domain.User) (domain.User, domain.ResponseError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, domain.NewReponseError(
			domain.ErrorInternal, err.Error(),
		)
	}
	user.Password = string(hashedPassword)

	result := uRepo.conn.Create(&user)
	if result.Error != nil {
		if result.Error.Error() == "Error 1062 (23000): Duplicate entry 'admin@gmail.com' for key 'users.email'" {
			return domain.User{}, domain.NewReponseError(
				domain.ErrorAlreadyExists, fmt.Sprintf("Email %v has been used", user.Email),
			)
		}
		return domain.User{}, domain.NewReponseError(
			domain.ErrorInternal, result.Error.Error(),
		)
	}

	return user, nil
}

// func (uRepo *userRepository) CheckIfExistsByEmail(email string) (bool, domain.ResponseError) {
// 	var user domain.User
// 	result := uRepo.conn.Where("email = ?", email).Find(&user)
// 	if result.Error != nil {
// 		return false, domain.NewReponseError(
// 			domain.ErrorInternal, result.Error.Error(),
// 		)
// 	}
// 	return result.RowsAffected > 0, nil
// }
