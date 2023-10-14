package main

import (
	"Go-Get/download"
	"Go-Get/getname"
	"Go-Get/merge"
	"Go-Get/way"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()

	// 设置websocket
	// 初始化 WebSocket 处理器
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// var conn *websocket.Conn

	r.GET("/ws", func(c *gin.Context) {
		// 升级 HTTP 连接为 WebSocket 连接
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("升级为 WebSocket 失败:", err)
			return
		}

		// 处理 WebSocket 连接
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("读取消息失败:", err)
				break
			}
			if messageType == websocket.TextMessage {
				// 处理文本消息
				message := string(p)
				// 在这里处理收到的消息
				fmt.Println(message)
			}
		}
		// 你可以在这里编写处理 WebSocket 连接的逻辑
	})

	r.LoadHTMLGlob("templates/*.html")
	// 设置静态文件服务
	r.Static("/static", "./static")

	//设置默认页面
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/download", func(c *gin.Context) {
		downloadURL := c.Query("downloadUrl")
		url := downloadURL
		//通过HEAD方法获取文件信息
		resp, err := http.Head(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
			return
		}
		defer resp.Body.Close()

		// 获取Content-Disposition头部字段，里面可能包含文件名
		contentDisposition := resp.Header.Get("Content-Disposition")
		fileName := ""
		if contentDisposition != "" {
			fileName, err = getname.ExtractFileNameFromContentDisposition(contentDisposition)
			if err != nil {
				// 处理提取文件名失败的情况
				fmt.Println("Error extracting file name:", err)
			}
		}

		// 如果无法从Content-Disposition中提取到文件名，则使用URL中的文件名
		if fileName == "" {
			fileName = filepath.Base(url)
		}

		// 从head获取文件大小，传入需要解析的字符串、进制（10进制）、结果类型（int64）
		contentLength, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
		chunkSize := contentLength / 5 // 分片大小 分成5个分片，你可以根据需要更改分片数

		var wg sync.WaitGroup
		progressCh := make(chan int, 5) // 同时下载的分片数

		// 下载分片
		for i := int64(0); i < 5; i++ {
			wg.Add(1)
			go download.DownloadChunk(url, i, chunkSize, &wg, progressCh)
		}

		// 监听分片下载进度，更新总进度条
		var totalProgress int64
		go func() {
			for p := range progressCh {
				totalProgress += int64(p)
				fmt.Printf("Total Progress: %.2f%%\n", float64(totalProgress)/float64(contentLength)*100)
			}
		}()

		// 等待所有分片下载完成
		wg.Wait()
		close(progressCh)

		// 合并分片文件
		merge.MergeChunks(5, "/Users/siky/Desktop/"+fileName)

		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
	})

	r.GET("/torrent", func(c *gin.Context) {
		// 种子文件路径
		torrentFilePath := "/Users/siky/go/src/Go-Get/test.torrent"

		// 传入种子文件路径和下载目录
		download.DownloadTorrentFile(torrentFilePath, "/Users/siky/go/src/Go-Get")
		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
		//删除 .torrent.db 文件
		dbFilePath := filepath.Join("/Users/siky/go/src/Go-Get", ".torrent.db")
		err := os.Remove(dbFilePath)
		if err != nil {
			log.Println("删除 .torrent.db 文件时出错:", err)
		} else {
			log.Println(".torrent.db 文件已成功删除")
		}
	})

	r.GET("/magnet", func(c *gin.Context) {
		// 磁力链接
		magnetURL := c.Query("magnetURL")
		//url := "magnet:?
		// 传入磁力链接和下载目录
		outputCh := make(chan string, 100)
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
		// 启动Goroutine将输出信息发送到WebSocket客户端
		// go func() {
		// 	for {
		// 		output := <-outputCh
		// 		err := conn.WriteMessage(websocket.TextMessage, []byte(output))
		// 		if err != nil {
		// 			way.SendOutput(outputCh, "发送消息失败:%v", err)
		// 			break
		// 		}
		// 	}
		// }()
	})

	r.Run(":9000")
}
