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
func DownloadChunk(url string, chunkIndex, chunkSize int64, wg *sync.WaitGroup, progressCh chan<- int64) {
	//传入URL、分片的索引、分片的大小、WaitGroup、发送分片下载进度的管道

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
	fileName := fmt.Sprintf("chunk_%d.tmp", chunkIndex) //分片文件名，格式为chunk_0.tmp、chunk_1.tmp等
	file, err := os.Create(fileName)                    //创建文件，返回文件指针
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

	// 将分块下载的数据从 HTTP 响应体复制到分块文件，并通过进度条实时显示下载进度。https://github.com/schollz/progressbar配套的
	//从resp.Body复制到file，同时也复制到bar，这样就可以实时显示下载进度了
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		fmt.Println("Error copying chunk:", err)
		return
	}

	// 下载完成时，发送分片下载进度给管道
	progressCh <- resp.ContentLength
}
