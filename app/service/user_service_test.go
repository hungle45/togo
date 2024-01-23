package service_test

import (
	"testing"
	"togo/app/service"
	"togo/domain"
	"togo/domain/mock"
	"togo/utils"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func createRandomUser() domain.User {
	return domain.User{
		Email:    utils.RandomEmail(),
		Password: utils.RandomPassword(),
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockUserRepo)

	t.Run("Test case 1: Success", func(t *testing.T) {
		user := createRandomUser()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		require.NoError(t, err)
		mockUserRepo.EXPECT().
			GetUserByEmail(user.Email).
			Return(domain.User{
				Email:    user.Email,
				Password: string(hashedPassword),
			}, nil)

		token, rerr := userService.Login(user)
		require.Nil(t, rerr)
		require.NotEmpty(t, token)
	})

	t.Run("Test case 2: Error (password not correct)", func(t *testing.T) {
		user := createRandomUser()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(utils.RandomPassword()), bcrypt.DefaultCost)
		require.NoError(t, err)
		mockUserRepo.EXPECT().
			GetUserByEmail(user.Email).
			Return(domain.User{
				Email:    user.Email,
				Password: string(hashedPassword),
			}, nil)

		token, rerr := userService.Login(user)
		require.Equal(t, rerr.ErrorType(), domain.ErrorUnauthenticated)
		require.Equal(t, token, "")
	})

	t.Run("Test case 3: Error (user not found)", func(t *testing.T) {
		user := createRandomUser()
		mockUserRepo.EXPECT().
			GetUserByEmail(user.Email).
			Return(domain.User{}, domain.NewReponseError(domain.ErrorNotFound, ""))

		token, rerr := userService.Login(user)
		require.Equal(t, rerr.ErrorType(), domain.ErrorUnauthenticated)
		require.Equal(t, token, "")
	})
}

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockUserRepo)

	t.Run("Test case 1: Success", func(t *testing.T) {
		user := createRandomUser()
		mockUserRepo.EXPECT().
			CreateUser(user).
			Return(domain.User{}, nil)

		rerr := userService.SignUp(user)
		require.Nil(t, rerr)
	})

	t.Run("Test case 2: Error (internal)", func(t *testing.T) {
		user := createRandomUser()
		mockUserRepo.EXPECT().
			CreateUser(user).
			Return(domain.User{}, domain.NewReponseError(domain.ErrorInternal, ""))

		rerr := userService.SignUp(user)
		require.Equal(t, rerr.ErrorType(), domain.ErrorInternal)
	})
}
