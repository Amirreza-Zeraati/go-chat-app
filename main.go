package main

import (
	"github.com/gin-gonic/gin"
	"go-chat-app/controllers/chatController"
	"go-chat-app/initializers"
	middleware "go-chat-app/middleware/auth"
	"go-chat-app/routes"
)

func init() {
	initializers.LoadEnvFile()
	initializers.ConnectToDB()
	initializers.Migrate()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	routes.AuthRoutes(r)
	r.GET("/chat", middleware.RequireAuth, chatController.ChatHub)

	err := r.Run()
	if err != nil {
		return
	}
}
