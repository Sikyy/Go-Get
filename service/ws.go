package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 在全局范围内定义一个 map 用于跟踪已连接的 WebSocket 客户端
var connectedClients = make(map[*websocket.Conn]bool)

// 设置websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Ws(c *gin.Context) {
	// 使用WebSocket处理程序
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// 添加连接到已连接客户端 Map
	connectedClients[ws] = true

	// 处理 WebSocket 消息
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// 处理断开连接
			log.Println("连接失败:", err)
			delete(connectedClients, ws)
			return
		}
	}
}
