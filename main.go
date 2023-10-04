package main

import (
	"Go-Get/download"
	"Go-Get/getname"
	"Go-Get/merge"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()

	r.GET("/download", func(c *gin.Context) {
		url := "https://cn-gdfs-ct-01-11.bilivideo.com/upgcxcode/47/51/1220655147/1220655147-1-16.mp4?e=ig8euxZM2rNcNbRVhwdVhwdlhWdVhwdVhoNvNC8BqJIzNbfq9rVEuxTEnE8L5F6VnEsSTx0vkX8fqJeYTj_lta53NCM=&uipk=5&nbs=1&deadline=1696318756&gen=playurlv2&os=bcache&oi=17621919&trid=00006b144ea1d7114a5eb0bd4e9c7b85de0dh&mid=0&platform=html5&upsig=d02b6da498b011fb53adcd4b139e1186&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,mid,platform&cdnid=60911&bvc=vod&nettype=0&f=h_0_0&bw=27530&logo=80000000"
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
		magnetURL := "magnet:?xt=urn:btih:4LREJASYZ6PGMSQODLXBSBFCWSGXQ2OR&dn=&tr=http%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=udp%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.prq.to%2Fannounce&tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce&tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud&tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce&tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce&tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce&tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce&tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce&tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce&tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce&tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce"

		// 传入磁力链接和下载目录
		download.DownloadMagnetFile(magnetURL, "/Users/siky/go/src/Go-Get")
		c.JSON(http.StatusOK, gin.H{"message": "Download completed"})
		// 删除 .torrent.db 文件
		dbFilePath := filepath.Join("/Users/siky/go/src/Go-Get", ".torrent.db")
		err := os.Remove(dbFilePath)
		if err != nil {
			log.Println("删除 .torrent.db 文件时出错:", err)
		} else {
			log.Println(".torrent.db 文件已成功删除")
		}
	})

	r.Run(":9000")
}
