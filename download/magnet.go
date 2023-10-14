package download

import (
	"Go-Get/getname"
	"Go-Get/way"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

type File struct {
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"is_dir"`
}

// DownloadMagnetFile 添加一个磁力链接并启动下载任务
func DownloadMagnetFile(magnetURL, downloadDir string, outputCh chan<- string) {
	// 创建一个新的 Torrent 客户端
	client, err := torrent.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}
	way.SendOutput(outputCh, "客户端初始化成功")
	defer client.Close()

	// 解析磁力链接，返回一个 TorrentFile 的指针
	torrentFile, err := client.AddMagnet(magnetURL)
	if err != nil {
		way.SendOutput(outputCh, "解析磁力链接失败:%v", err)
		log.Fatal(err)
	}
	//q：为什么会返回空指针
	//a：因为这个磁力链接是无效的，或者说这个磁力链接没有对应的种子文件

	// 等待解析磁力链接，当元数据下载完成后，该通道会被关闭，此时可以访问种子信息
	// 如果超时，就会执行 time.After() 中的代码
	select {
	case <-torrentFile.GotInfo():
		way.SendOutput(outputCh, "解析磁力链接成功")
		// 此处添加访问种子信息的代码
	case <-time.After(15 * time.Second):
		way.SendOutput(outputCh, "解析超时，建议检查一下网络情况，或是磁力链接是否失效")
		return // 或者执行其他超时后的操作
	}

	info := torrentFile.Info()

	way.SendOutput(outputCh, "总文件名称:%v", info.Name)
	way.SendOutput(outputCh, "总文件数量:%v", info.TotalLength())
	way.SendOutput(outputCh, "总文件大小:%v", len(info.Files))

	files := torrentFile.Files()
	for _, file := range files {
		way.SendOutput(outputCh, "文件名称:%v", file.Path())
		way.SendOutput(outputCh, "文件大小:%v", file.Length())
	}

	// 获取种子文件的原始文件名
	torrentFileName := filepath.Base(torrentFile.Info().Name)
	way.SendOutput(outputCh, "文件原始名称:%v", torrentFileName)
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

	torrentFile.AllowDataDownload() // 允许下载数据
	way.SendOutput(outputCh, "允许下载数据")

	// 下载所有文件
	torrentFile.DownloadAll()
	way.SendOutput(outputCh, "开始下载")

	// 等待下载完成
	client.WaitAll()
	way.SendOutput(outputCh, "下载完成:"+torrentFileName)
}
