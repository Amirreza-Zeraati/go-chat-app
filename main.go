package main

import (
	"github.com/gin-gonic/gin"
	"go-chat-app/controllers/chatController"
	"go-chat-app/initializers"
	middleware "go-chat-app/middleware/auth"
	"go-chat-app/routes"
	"html/template"
)

func init() {
	initializers.LoadEnvFile()
	initializers.ConnectToDB()
	initializers.Migrate()
}

func main() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"nl2br": chatController.Nl2br,
	})
	r.LoadHTMLGlob("templates/*")
	routes.AuthRoutes(r)
	r.GET("/chat", middleware.RequireAuth, func(c *gin.Context) {
		chatController.HandleChatPage(c)
	})
	r.GET("/ws", middleware.RequireAuth, func(c *gin.Context) {
		chatController.HandleConnections(c)
	})

	go chatController.HandleMessages()

	err := r.Run()
	if err != nil {
		return
	}
}
