package routes

import (
	"github.com/gin-gonic/gin"
	"go-chat-app/controllers/userController"
)

func AuthRoutes(r *gin.Engine) {
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)
}
