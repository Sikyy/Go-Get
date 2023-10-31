package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UploadTorrentInfo struct {
	Name         string    `json:"name"`          // 种子文件名称
	FilesNum     int       `json:"files_num"`     // 种子文件数量
	TotalLength  int64     `json:"total_length"`  // 种子文件大小
	InfoHash     string    `json:"info_hash"`     // 种子文件哈希
	InfoBytes    string    `json:"info_bytes"`    // 种子元信息数组
	Announce     string    `json:"announce"`      // 种子tracker地址
	Comment      string    `json:"comment"`       // 种子注释
	CreatedBy    string    `json:"created_by"`    // 种子创建者
	CreationDate int64     `json:"creation_date"` // 种子创建时间
	UploadTime   time.Time `json:"upload_time"`   // 种子上传时间
}

// 把uploadInfo发送到/uploadToMongoDB
func SendUploadInfoToMongoDB(uploadInfo UploadTorrentInfo) error {
	// 将 uploadInfo 结构体编码为 JSON 数据
	jsonData, err := json.Marshal(uploadInfo)
	if err != nil {
		return err
	}

	// 发送 HTTP POST 请求将 JSON 数据发送到/uploadToMongoDB路由
	url := "http://localhost:9000/uploadToMongoDB"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// 上传成功
		fmt.Println("UploadInfo sent to MongoDB successfully")
	} else {
		fmt.Println("UploadInfo sending to MongoDB failed")
	}

	return nil
}

func SendUploadInfoToTestMongoDB(uploadInfo UploadTorrentInfo) error {
	// 将 uploadInfo 结构体编码为 JSON 数据
	jsonData, err := json.Marshal(uploadInfo)
	if err != nil {
		return err
	}

	// 发送 HTTP POST 请求将 JSON 数据发送到/uploadToMongoDB路由
	url := "http://localhost:9000/testUploadToMongoDB"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// 上传成功
		fmt.Println("UploadInfo sent to MongoDB successfully")
	} else {
		fmt.Println("UploadInfo sending to MongoDB failed")
	}

	return nil
}
