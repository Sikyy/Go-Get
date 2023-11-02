package service

import (
	"Go-Get/data"
	"Go-Get/download"
	"Go-Get/way"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TorrentDownload(c *gin.Context) {
	// 种子文件路径
	//torrentFilePath := "/Users/siky/go/src/Go-Get/test.torrent"
	torrentFilePath := c.Query("torrentFilePath")
	// torrentFilePath := "/Users/siky/go/src/Go-Get/uploads/"
	outputCh := make(chan string, 10000)
	// 传入种子文件路径和下载目录
	go func() {
		defer close(outputCh)
		uploadinfo := download.DownloadTorrentFile(torrentFilePath, "/Users/siky/go/src/Go-Get", outputCh)

		// 在这里将 uploadInfo 发送到/uploadToMongoDB路由
		err := data.SendUploadInfoToMongoDB(uploadinfo)
		if err != nil {
			way.SendOutput(outputCh, "上传种子信息到MongoDB时出错:%v", err)
		} else {
			way.SendOutput(outputCh, "种子信息已成功上传到MongoDB")
		}
		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})

		//删除 .torrent.db 文件
		dbFilePath := filepath.Join("/Users/siky/go/src/Go-Get", ".torrent.db")
		err = os.Remove(dbFilePath)
		if err != nil {
			way.SendOutput(outputCh, "删除 .torrent.db 文件时出错:%v", err)
		} else {
			way.SendOutput(outputCh, ".torrent.db 文件已成功删除")
		}
		way.SendOutput(outputCh, "------------------------该文件下载完成------------------------")
		// 在这里将 uploadInfo 发送到/uploadToMongoDB路由
	}()
	// 处理输出消息并发送到 WebSocket 客户端
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
