package chatController

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat-app/initializers"
	"go-chat-app/models"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func Nl2br(text string) template.HTML {
	safe := template.HTMLEscapeString(text) // escape HTML
	safe = strings.ReplaceAll(safe, "\n", "<br>")
	return template.HTML(safe) // mark safe for template
}

func HandleChatPage(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)
	var chats []models.Chat
	initializers.DB.Order("created_at asc").Find(&chats)
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"username": user.Name,
		"chats":    chats,
	})
}

func HandleConnections(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Println("Close error:", err)
		}
	}(ws)
	clients[ws] = true
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
		broadcast <- string(msg)
		var chat models.Chat
		if err := json.Unmarshal(msg, &chat); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}
		chat = models.Chat{UserID: user.ID, Name: user.Name, Text: chat.Text}
		result := initializers.DB.Create(&chat)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create chat: " + result.Error.Error(),
			})
			return
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				err := client.Close()
				if err != nil {
					return
				}
				delete(clients, client)
			}
		}
	}
}
