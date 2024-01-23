package http

import (
	"net/http"
	http_utils "togo/app/delivery/http/utils"
	"togo/domain"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type userForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type userHTTPHandler struct {
	userService domain.UserService
}

func NewUserHTTPHandler(r *gin.RouterGroup, userService domain.UserService) {
	handler := userHTTPHandler{userService: userService}

	userRouter := r.Group("/users")
	{
		userRouter.POST("/login", handler.login)
		userRouter.POST("/signup", handler.signUp)
	}
}

func (handler *userHTTPHandler) login(c *gin.Context) {
	form, err := validateUserForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error()))
		return
	}

	user := domain.User{
		Email:    form.Email,
		Password: form.Password,
	}

	token, rerr := handler.userService.Login(user)
	if rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, rerr.Message()))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithData(
		http_utils.ReponseStatusSuccess,
		map[string]interface{}{"token": token}),
	)
}

func (handler *userHTTPHandler) signUp(c *gin.Context) {
	form, err := validateUserForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, err.Error()))
		return
	}

	user := domain.User{
		Email:    form.Email,
		Password: form.Password,
	}
	if rerr := handler.userService.SignUp(user); rerr != nil {
		c.JSON(http_utils.GetStatusCode(rerr), http_utils.ResponseWithMessage(
			http_utils.ResponseStatusFail, rerr.Message()))
		return
	}

	c.JSON(http.StatusOK, http_utils.ResponseWithMessage(
		http_utils.ReponseStatusSuccess, "account has been created"))
}

func validateUserForm(c *gin.Context) (userForm, error) {
	form := userForm{}
	if err := c.ShouldBindBodyWith(&form, binding.JSON); err != nil {
		return userForm{}, err
	}
	return form, nil
}
