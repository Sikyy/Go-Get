package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"siky-idm/download"
	"siky-idm/getname"
	"siky-idm/merge"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/download", func(c *gin.Context) {
		url := "https://cn-gddg-ct-01-12.bilivideo.com/upgcxcode/48/28/1282432848/1282432848-1-16.mp4?e=ig8euxZM2rNcNbRVhwdVhwdlhWdVhwdVhoNvNC8BqJIzNbfq9rVEuxTEnE8L5F6VnEsSTx0vkX8fqJeYTj_lta53NCM=&uipk=5&nbs=1&deadline=1696264842&gen=playurlv2&os=bcache&oi=17627301&trid=00008c2b90aa305d4b34b771bdb3db3b7b48h&mid=0&platform=html5&upsig=9d53d11b621f908fdc61ae6fd0339e2d&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,mid,platform&cdnid=61312&bvc=vod&nettype=0&f=h_0_0&bw=51461&logo=80000000"
		resp, err := http.Head(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
			return
		}
		defer resp.Body.Close()

		// 获取Content-Disposition头部字段
		contentDisposition := resp.Header.Get("Content-Disposition")
		var fileName string

		if contentDisposition != "" {
			var err error
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
		merge.MergeChunks(5, "/Users/siky/Desktop/"+fileName)

		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
	})

	r.Run(":9000")
}

func extractFileNameFromContentDisposition(contentDisposition string) {
	panic("unimplemented")
}
