package service

import (
	"Go-Get/data"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TestTotestUploadToMongoDB(c *gin.Context) {
	// 测试上传数据给testUploadToMongoDB路由
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
}

func TestUploadToMongoDB(c *gin.Context) {
	var receivedData data.UploadTorrentInfo
	if err := c.ShouldBindJSON(&receivedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 在这里处理接收到的数据
	fmt.Println("Received data:", receivedData)

	c.JSON(http.StatusOK, receivedData)
}
