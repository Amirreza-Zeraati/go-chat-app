package chatController

import (
	"github.com/gin-gonic/gin"
	"go-chat-app/models"
	"net/http"
)

func ChatHub(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"user": user,
	})
}
