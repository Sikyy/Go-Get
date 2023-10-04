package download

import (
	"Go-Get/getname"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

// DownloadMagnetFile 添加一个磁力链接并启动下载任务
func DownloadMagnetFile(magnetURL, downloadDir string) {
	// 创建一个新的 Torrent 客户端
	client, err := torrent.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("初始化客户端成功")
	defer client.Close()

	// 解析磁力链接，返回一个 TorrentFile 的指针
	torrentFile, err := client.AddMagnet(magnetURL)
	if err != nil {
		fmt.Println("解析磁力链接时出错:", err)
		log.Fatal(err)
	}
	//q：为什么会返回空指针
	//a：因为这个磁力链接是无效的，或者说这个磁力链接没有对应的种子文件

	// 等待解析磁力链接，当元数据下载完成后，该通道会被关闭，此时可以访问种子信息
	// 如果超时，就会执行 time.After() 中的代码
	select {
	case <-torrentFile.GotInfo():
		fmt.Println("解析磁力链接成功")
		// 此处添加访问种子信息的代码
	case <-time.After(15 * time.Second):
		fmt.Println("解析超时，建议检查一下网络情况，或是磁力链接是否失效")
		return // 或者执行其他超时后的操作
	}

	info := torrentFile.Info()

	log.Println("Torrent Name:", info.Name)
	log.Println("Number of Files:", len(info.Files))
	log.Println("Total Size:", info.TotalLength())

	files := torrentFile.Files()
	for _, file := range files {
		log.Println("File Name:", file.Path())
		log.Println("File Size:", file.Length())
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
	fmt.Println("Save path:", savePath)

	// 创建目录（如果不存在），默认权限是0777
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		log.Fatal(err)
		fmt.Println("创建目录失败")
	} else {
		fmt.Println("创建目录成功")
	}

	// 设置下载目录，如果是部分下载，就可以指定保存位置，全部下载保存位置就是默认的下载目录
	downloadDir = savePath

	// 更改当前工作目录到下载目录
	if err := os.Chdir(savePath); err != nil {
		log.Fatal(err)
	}

	torrentFile.AllowDataDownload() // 允许下载数据
	fmt.Println("允许下载数据")

	// 下载所有文件
	torrentFile.DownloadAll()
	fmt.Println("开始下载")

	// 等待下载完成
	client.WaitAll()
	log.Println("Torrent downloaded:", torrentFile.Info().Name)
}
