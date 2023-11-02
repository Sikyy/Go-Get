package router

import (
	"Go-Get/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 设置跨域访问配置
	r.Use(cors.Default())

	r.LoadHTMLGlob("templates/*.html")
	// 设置静态文件服务
	r.Static("/static", "./static")

	//websocket连接相关
	r.GET("/ws", service.Ws)

	//页面相关
	r.GET("/", service.HomePage)
	r.GET("/background", service.BackgroundPage)

	//下载相关
	//HTTP/HTTPS下载
	r.GET("/download", service.HttpDownload)
	//Torrent下载
	r.GET("/torrent", service.TorrentDownload)
	//Magent下载
	r.GET("/magnet", service.MagnetDownload)

	//上传相关
	r.POST("/upload", service.Upload)
	r.POST("/uploadToMongoDB", service.UploadToMongoDB)

	//测试相关
	//测试上传数据给testUploadToMongoDB路由
	r.GET("/test", service.TestTotestUploadToMongoDB)
	//测试上传数据给MongoDB
	r.POST("/testUploadToMongoDB", service.TestUploadToMongoDB)
	return r
}
