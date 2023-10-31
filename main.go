package main

import (
	"Go-Get/data"
	"Go-Get/download"
	"Go-Get/getname"
	"Go-Get/merge"
	"Go-Get/way"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()

	// 全局变量用于存储上传的 torrent 文件路径
	// var uploadedTorrentPath string

	// 在全局范围内定义一个 map 用于跟踪已连接的 WebSocket 客户端
	var connectedClients = make(map[*websocket.Conn]bool)

	// 设置websocket
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	r.GET("/ws", func(c *gin.Context) {
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

		// 模拟实时输出内容
		// for i := 0; i < 5; i++ {
		// 	message := "This is content number " + fmt.Sprint(i+1)
		// 	//睡眠1秒
		// 	time.Sleep(time.Second)
		// 	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		// 		return
		// 	}
		// }
	})

	r.LoadHTMLGlob("templates/*.html")
	// 设置静态文件服务
	r.Static("/static", "./static")

	// 设置跨域访问配置
	r.Use(cors.Default())

	// 建立MongoDB连接
	client, err := data.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

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

	r.POST("/upload", func(c *gin.Context) {
		// 从请求中读取文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 将文件保存到指定目录
		uploadDir := "/Users/siky/go/src/Go-Get/static/uploads"
		//如果文件夹不存在则创建文件夹
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("文件夹创建成功")

		uploadedFilePath := filepath.Join(uploadDir, file.Filename)

		// 保存文件到指定路径
		if err := c.SaveUploadedFile(file, uploadedFilePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "文件上传成功", "uploadedFilePath": uploadedFilePath})

	})

	r.POST("/uploadToMongoDB", func(c *gin.Context) {

		// 创建一个 downloadInfo 变量，用于存储 JSON 数据
		var downloadInfo data.UploadTorrentInfo

		// 从请求中读取 JSON 数据并解码到 downloadInfo 变量
		if err := c.BindJSON(&downloadInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//把json转化为BSON
		uploadinfo := bson.M{
			"name":          downloadInfo.Name,
			"files_num":     downloadInfo.FilesNum,
			"total_length":  downloadInfo.TotalLength,
			"info_hash":     downloadInfo.InfoHash,
			"info_bytes":    downloadInfo.InfoBytes,
			"announce":      downloadInfo.Announce,
			"comment":       downloadInfo.Comment,
			"created_by":    downloadInfo.CreatedBy,
			"creation_date": downloadInfo.CreationDate,
			"upload_time":   downloadInfo.UploadTime,
		}

		// 在此处将下载信息插入 MongoDB 数据库
		data.InsertDocument(client, uploadinfo, "Go-Get-MongoDB", "Torrents")

		c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
	})

	r.GET("/torrent", func(c *gin.Context) {
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
	})

	r.GET("/magnet", func(c *gin.Context) {
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
	})

	r.GET("/test", func(c *gin.Context) {
		uploadInfo := data.UploadTorrentInfo{
			Name:         "test",
			FilesNum:     1,
			TotalLength:  1,
			InfoHash:     "test",
			InfoBytes:    "test",
			Announce:     "test",
			Comment:      "test",
			CreatedBy:    "test",
			CreationDate: 1,
			UploadTime:   time.Now(),
		}

		if err := data.SendUploadInfoToMongoDB(uploadInfo); err != nil {
			fmt.Println("Error sending data:", err)
		}
	})

	r.POST("/testUploadToMongoDB", func(c *gin.Context) {
		var receivedData data.UploadTorrentInfo
		if err := c.ShouldBindJSON(&receivedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 在这里处理接收到的数据
		fmt.Println("Received data:", receivedData)

		c.JSON(http.StatusOK, receivedData)
	})

	r.Run(":9000")
}
