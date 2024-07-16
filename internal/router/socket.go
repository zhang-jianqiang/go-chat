package router

import (
	"chat-room/internal/server"
	"chat-room/pkg/global/log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunSocekt(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		return
	}
	log.Logger.Info("newUser", zap.String("newUser", user))
	// todo: 这里做认证校验
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &server.Client{
		Name:  user,
		Conn:  ws,
		Send:  make(chan []byte),
		Mutex: &sync.Mutex{},
	}

	server.MyServer.Register <- client
	go client.Read()
	go client.Write()
}
