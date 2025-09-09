package chatController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-chat-app/models"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func HandleChatPage(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)
	c.HTML(http.StatusOK, "chat.html", gin.H{
		"user": user,
	})
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {

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
