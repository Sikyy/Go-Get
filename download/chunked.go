package download

//文件分块下载

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// DownloadChunk 下载分片的函数
func DownloadChunk(url string, chunkIndex, chunkSize int64, wg *sync.WaitGroup, progressCh chan<- int) {
	//传入URL、分片的索引、大小、WaitGroup、进度通道

	defer wg.Done()

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// 计算分片的范围
	start := chunkIndex * chunkSize
	end := start + chunkSize - 1
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	// 发送请求并获取响应
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error executing HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// 创建分片文件
	fileName := fmt.Sprintf("chunk_%d.tmp", chunkIndex)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating chunk file:", err)
		return
	}
	defer file.Close()

	// 创建进度条
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("Chunk %d", chunkIndex),
	)

	// 复制响应体到分片文件，并同时更新进度条
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		fmt.Println("Error copying chunk:", err)
		return
	}

	// 发送分片下载进度
	progressCh <- int(resp.ContentLength)
}
