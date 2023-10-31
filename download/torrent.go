package download

import (
	"Go-Get/data"
	"Go-Get/getname"
	"Go-Get/way"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

// DownloadTorrentFile 添加一个种子文件并启动下载任务
func DownloadTorrentFile(torrentFilePath, downloadDir string, outputCh chan<- string) data.UploadTorrentInfo {
	// 创建一个新的Torrent客户端
	client, err := torrent.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}
	way.SendOutput(outputCh, "客户端初始化成功")
	defer client.Close()

	// 解析上传的种子文件，返回一个TorrentFile的指针
	torrentFile, err := client.AddTorrentFromFile(torrentFilePath)
	if err != nil {
		log.Fatal(err)
	}
	way.SendOutput(outputCh, "解析种子文件成功")

	// 创建一个 UploadTorrentInfo 变量并填充数据
	uploadInfo := data.UploadTorrentInfo{
		Name:         torrentFile.Info().Name,
		FilesNum:     len(torrentFile.Info().Files) + 1,
		TotalLength:  torrentFile.Info().TotalLength(),
		InfoHash:     torrentFile.InfoHash().HexString(),
		InfoBytes:    string(torrentFile.Metainfo().InfoBytes),
		Announce:     torrentFile.Metainfo().Announce,
		Comment:      torrentFile.Metainfo().Comment,
		CreatedBy:    torrentFile.Metainfo().CreatedBy,
		CreationDate: torrentFile.Metainfo().CreationDate,
		UploadTime:   time.Now(),
	}

	// 访问解析后的种子信息,给前端返回种子信息
	way.SendOutput(outputCh, "种子文件名称:%v", torrentFile.Info().Name)
	way.SendOutput(outputCh, "种子文件数量:%v", len(torrentFile.Info().Files)+1)
	way.SendOutput(outputCh, "种子文件大小:%v", torrentFile.Info().TotalLength())

	// 获取 Tracker 列表
	trackers := torrentFile.Metainfo().AnnounceList
	for _, tracker := range trackers {
		way.SendOutput(outputCh, "Tracker URL:%v", tracker)
	}

	files := torrentFile.Files()
	for _, file := range files {
		// fmt.Println("File Name:", file.Path())
		// fmt.Println("File Size:", file.Length())
		way.SendOutput(outputCh, "文件名称:%v", file.Path())
		way.SendOutput(outputCh, "文件大小:%v", file.Length())
	}

	// 获取种子文件的原始文件名
	torrentFileName := filepath.Base(torrentFile.Info().Name)
	fmt.Println("Torrent file name:", torrentFileName)
	// 获取种子文件的原始文件名（不包含扩展名）
	torrentFileNameWithoutExtension := strings.TrimSuffix(torrentFileName, filepath.Ext(torrentFileName))
	//使用正则表达式匹配动画名称
	animeName := getname.ExtractAnimeName(torrentFileNameWithoutExtension)

	// 设置保存目录
	savePath := filepath.Join(downloadDir, animeName)
	way.SendOutput(outputCh, "保存目录:%v", savePath)

	// 创建目录（如果不存在），默认权限是0777
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		log.Fatal(err)
		way.SendOutput(outputCh, "创建目录失败:%v", err)
	} else {
		way.SendOutput(outputCh, "创建目录成功")
	}

	// 设置下载目录，如果是部分下载，就可以指定保存位置，全部下载保存位置就是默认的下载目录
	downloadDir = savePath

	// 更改当前工作目录到下载目录
	if err := os.Chdir(savePath); err != nil {
		log.Fatal(err)
	}
	// 下载所有文件或者指定文件
	// 选择要下载的文件（示例中选择第一个文件）
	// selectedFile := files[0]
	torrentFile.AllowDataDownload() // 允许下载数据
	way.SendOutput(outputCh, "允许下载数据")

	torrentFile.DownloadAll() //开始下载
	way.SendOutput(outputCh, "开始下载")

	// 开始下载
	torrentFile.DownloadAll()

	// 等待下载完成
	client.WaitAll()
	way.SendOutput(outputCh, "下载完成:"+torrentFileName)

	// 在这里将 uploadInfo 发送到/uploadToMongoDB路由
	// err = SendUploadInfoToMongoDB(uploadInfo)
	// if err != nil {
	// 	fmt.Println("Error sending uploadInfo to MongoDB")
	// 	log.Fatal(err)
	// }
	// way.SendOutput(outputCh, "上传种子信息到MongoDB成功")
	return uploadInfo
}
