package repository_test

import (
	"testing"
	"time"
	"togo/domain"
	"togo/utils"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func createRandomUser(t *testing.T) domain.User {
	randomUser := domain.User{
		Email:    utils.RandomEmail(),
		Password: utils.RandomPassword(),
	}

	user, rerr := userRepository.CreateUser(randomUser)
	require.Empty(t, rerr)
	require.Equal(t, randomUser.Email, user.Email)
	require.True(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(randomUser.Password)) == nil)
	require.NotEmpty(t, user.CreatedAt.In(time.UTC))
	require.NotEmpty(t, user.UpdatedAt.In(time.UTC))

	user.Password = randomUser.Password
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	testUser := createRandomUser(t)

	user, rerr := userRepository.GetUserByID(testUser.ID)
	require.Empty(t, rerr)
	require.Equal(t, testUser.ID, user.ID)
	require.Equal(t, testUser.Email, user.Email)
	require.True(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testUser.Password)) == nil)
	require.Equal(t, testUser.CreatedAt.In(time.UTC), user.CreatedAt.In(time.UTC))
	require.Equal(t, testUser.UpdatedAt.In(time.UTC), user.UpdatedAt.In(time.UTC))
}

func TestGetUserByEmail(t *testing.T) {
	testUser := createRandomUser(t)

	user, rerr := userRepository.GetUserByEmail(testUser.Email)
	require.Empty(t, rerr)
	require.Equal(t, testUser.ID, user.ID)
	require.Equal(t, testUser.Email, user.Email)
	require.True(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testUser.Password)) == nil)
	require.Equal(t, testUser.CreatedAt.In(time.UTC).In(time.UTC), user.CreatedAt.In(time.UTC).In(time.UTC))
	require.Equal(t, testUser.UpdatedAt.In(time.UTC).In(time.UTC), user.UpdatedAt.In(time.UTC).In(time.UTC))
}
