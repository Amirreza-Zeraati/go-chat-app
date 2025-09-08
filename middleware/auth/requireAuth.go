package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-chat-app/initializers"
	"go-chat-app/models"
	"net/http"
	"os"
	"time"
)

func RequireAuth(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	secret := []byte(os.Getenv("SECRET"))
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || token == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["exp"], claims["sub"])
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
