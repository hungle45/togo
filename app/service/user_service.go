package usecase

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"togo/config"
	"togo/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{userRepo: userRepo}
}

func (uS *userService) createToken(user domain.User) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(secretKey)
}

func (uS *userService) extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 || jwtToken[0] != "Bearer" {
		return "", errors.New("invalid formatted authorization header")
	}

	return jwtToken[1], nil
}

func (uS *userService) ParseToken(header string) (uint, domain.ResponseError) {
	jwtToken, err := uS.extractBearerToken(header)
	if err != nil {
		return 0, domain.NewReponseError(domain.ErrorUnauthenticated, err.Error())
	}
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // Check if the signing method is HMAC
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return 0, domain.NewReponseError(domain.ErrorUnauthenticated, "invalid credentials")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, domain.NewReponseError(domain.ErrorUnauthenticated, "invalid credentials")
	}

	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		return 0, domain.NewReponseError(domain.ErrorUnauthenticated, "invalid credentials")
	}

	userID := uint(claims["sub"].(float64))
	return userID, nil
}

func (uS *userService) Login(user domain.User) (string, domain.ResponseError) {
	res, rerr := uS.userRepo.GetUserByEmail(user.Email)
	if rerr != nil || bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password)) != nil {
		return "", domain.NewReponseError(domain.ErrorUnauthenticated, "invalid credentials")
	}

	token, err := uS.createToken(res)
	if err != nil {
		return "", domain.NewReponseError(domain.ErrorInternal, "unable to generate token")
	}

	return token, nil
}

func (uS *userService) SignUp(user domain.User) domain.ResponseError {
	if _, rerr := uS.userRepo.CreateUser(user); rerr != nil {
		return rerr
	}

	return nil
}

func (uS *userService) CreateAdmin(cfg *config.Config) {
	_, rerr := uS.userRepo.GetUserByEmail(cfg.App.Server.Admin.Email)
	if rerr == nil {
		return
	}

	admin := domain.User{
		Email:    cfg.App.Server.Admin.Email,
		Password: cfg.App.Server.Admin.Password,
		Role:     domain.AdminRole,
	}

	_, rerr = uS.userRepo.CreateUser(admin)
	if rerr != nil {
		log.Fatal("Error creating admin user")
	}
}
