package main

import (
	"siky-idm/download"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个 Gin 引擎
	r := gin.Default()

	// 定义一个路由处理程序来触发下载功能
	r.GET("/download", func(c *gin.Context) {
		// 调用你的下载功能函数
		download.ChunkedDownload(c)
	})

	// 启动 HTTP 服务器并监听在 8080 端口
	r.Run(":8000")
}
