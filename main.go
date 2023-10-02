package main

import (
	"fmt"
	"net/http"
	"siky-idm/download"
	"siky-idm/merge"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/download", func(c *gin.Context) {
		url := "https://acg.rip/t/288868.torrent"
		resp, err := http.Head(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
			return
		}
		defer resp.Body.Close()

		contentLength, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
		chunkSize := contentLength / 5 // 分成5个分片，你可以根据需要更改分片数

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
		merge.MergeChunks(5, "merged_file.torrent")

		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
	})

	r.Run(":9000")
}
