package service

import (
	"Go-Get/download"
	"Go-Get/way"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func MagnetDownload(c *gin.Context) {
	// 磁力链接
	magnetURL := c.Query("magnetURL")
	//url := "magnet:?
	// 传入磁力链接和下载目录
	outputCh := make(chan string, 10000)
	go func() {
		defer close(outputCh)
		download.DownloadMagnetFile(magnetURL, "/Users/siky/go/src/Go-Get", outputCh)
		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
		// 删除 .torrent.db 文件
		dbFilePath := filepath.Join("/Users/siky/go/src/Go-Get", ".torrent.db")
		err := os.Remove(dbFilePath)
		if err != nil {
			way.SendOutput(outputCh, "删除 .torrent.db 文件时出错:%v", err)
		} else {
			way.SendOutput(outputCh, ".torrent.db 文件已成功删除")
		}
		way.SendOutput(outputCh, "------------------------该文件下载完成------------------------")
	}()
	//处理输出消息并发送到 WebSocket 客户端
	go func() {
		for message := range outputCh {
			for client := range connectedClients {
				if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					fmt.Println("Failed to send data to client:", err)
				}
			}
		}
	}()
}
