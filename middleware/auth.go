package middleware

import (
	"FlashSaleGo/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

func RequireAuth(service service.IUserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
			}

			return []byte(os.Getenv("SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//check expiration
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				//token expired
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
				c.Abort()
				return
			}
			user, err := service.GetUserByUsername(claims["sub"].(string))
			if err != nil || user == nil {
				//failed to get user
				c.Redirect(http.StatusTemporaryRedirect, "/user/login")
				c.Abort()
				return
			}
			c.Set("user", user)
			c.Next()
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			c.Abort()
			return
		}
	}
}
