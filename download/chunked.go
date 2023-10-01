package download

//文件分块下载

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/cheggaaa/pb"
	"github.com/gin-gonic/gin"
)

// ChunkedDownload 分块下载
func ChunkedDownload(c *gin.Context) {
	var copyfilepath string        //复制文件路径
	var storebreakpointpath string //存储断点信息文件路径
	// 获取要下载的文件路径
	downloadfilePath := "shell01.png"
	copyfolderPath := "/Users/siky/Desktop"
	// 创建文件夹，如果不存在的话
	err := os.MkdirAll(copyfolderPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}
	// 指定要保存的文件路径
	filePath := copyfolderPath + "/1.png"

	// 创建文件
	copyfile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer copyfile.Close()

	// 打开文件
	file, err := os.Open(downloadfilePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// 获取文件大小
	fileSize := fileInfo.Size()

	//统计分片数
	var scount int64                           //分片数
	desiredChunkSize := int64(1 * 1024 * 1024) // 每片分片大小1 MB（可以根据需要调整）
	scount = fileSize / desiredChunkSize

	// 如果文件大小不能整除 desiredChunkSize，增加一个分片以处理剩余部分
	if fileSize%desiredChunkSize != 0 {
		scount++
	}
	fmt.Printf("文件总大小：%v, 分片数：%v,每个分片大小：%v\n", fileSize, scount, desiredChunkSize)
	//创建下载文件的副本
	copyFile, err := os.OpenFile(copyfilepath, os.O_CREATE|os.O_WRONLY, os.ModePerm) //如果文件不存在则创建它
	if err != nil {
		fmt.Println(err)
	}
	defer copyFile.Close() //关闭下载副本文件

	//创建断点信息文件
	storgeFile, err := os.OpenFile(storebreakpointpath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	//断点信息文件需要到条件手动删除

	//设置当前分片的索引
	var currentindex int64 = 0 //当前分片索引
	wg := sync.WaitGroup{}     //创建同步等待组
	count := scount            //让count等于分片数
	var downloadedSize int64   // 用于跟踪已下载的总数据量
	fmt.Println("开始下载")
	totalBar := pb.StartNew(int(count))
	progressBar := 1                              //进度条标记点，每下载1%的数据，更新一次进度条
	for ; currentindex < scount; currentindex++ { //有多少个循环就开启多少个协程
		wg.Add(1)
		go func(current int64) { //定义一个匿名方法处理你需要协程处理的业务逻辑
			p := pb.New(int(desiredChunkSize)).Prefix(fmt.Sprintf("%dst", current+1))
			//fmt.Sprint((current+1))+"st" 是进度条的名称。它使用 fmt.Sprint() 函数将 current+1 转换为字符串，并加上 st 后缀。这样，进度条的名称将是 1st、2st、3st 等
			b := make([]byte, 1024)                              //数据缓冲区，每次读取1kb
			bs := make([]byte, 16)                               //断点信息缓冲区，每次读取16字节
			currentIndex, _ := storgeFile.ReadAt(bs, current*16) //读取断点信息，从断点信息文件 storgeFile 中读取数据，并将其存储在 bs 切片中
			//使用正则表达式从断点信息中提取当前分片的已下载大小（整数），防止出现小数点精度问题
			reg := regexp.MustCompile(`\d+`)
			countStr := reg.FindString(string(bs[:currentIndex])) //将 bs 切片转换为字符串，并使用正则表达式提取整数
			total, _ := strconv.ParseInt(countStr, 10, 0)         //将已下载大小转化成int64类型
			for {

				if total >= desiredChunkSize {
					wg.Done() //结束当前协程
					break
				}
				//从原始文件中读取数据
				realread, err := file.ReadAt(b, current*desiredChunkSize+total) //实际读取的字节数
				if err == io.EOF {
					wg.Done()
					break
				}
				//写入到下载副本文件中
				copyFile.WriteAt(b, current*desiredChunkSize+total)
				storgeFile.WriteAt([]byte(strconv.FormatInt(total, 10)+" "), current*16)

				// 更新已下载的总数据量
				atomic.AddInt64(&downloadedSize, int64(realread))

				// 更新总进度条
				totalBar.Set(int(downloadedSize))
				totalBar.Increment() // 增加总进度条的当前值

				//将已下载的数据大小（total 变量）以字符串形式写入到断点信息文件中，以便在下次继续下载时使用。
				total += int64(realread)

				if total >= desiredChunkSize/100*int64(progressBar) { //这个条件检查是否达到了下一个进度条标记点，如果是，则更新进度条
					progressBar++
					p.Add(int(desiredChunkSize / 100))
				}
			}

		}(currentindex)
	}
	wg.Wait() //协程等待
	storgeFile.Close()
	os.Remove(storebreakpointpath) //删除断点信息文件
	fmt.Println("下载完成")
	// 完成总进度条
	totalBar.Finish()
}
