package middleware

import (
	http_utils "togo/app/delivery/http/utils"
	"togo/domain"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(userService domain.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: Bearer <token>
		userID, rerr := userService.ParseToken(c.GetHeader("Authorization"))
		if rerr != nil {
			c.AbortWithStatusJSON(http_utils.GetStatusCode(rerr),
				http_utils.ResponseWithMessage(http_utils.ResponseStatusFail, "invalid token"))
		}
		c.Set("userID", userID)

		c.Next()
	}
}
