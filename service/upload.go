package service

import (
	"Go-Get/data"
	"Go-Get/dbinit"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 上传Torrent文件到服务器
func Upload(c *gin.Context) {
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

}

// 上传Torrent文件到MongoDB
func UploadToMongoDB(c *gin.Context) {
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
	data.InsertDocument(dbinit.Client, uploadinfo, "Go-Get-MongoDB", "Torrents")

	c.JSON(http.StatusOK, gin.H{"message": "Data inserted successfully"})
}
