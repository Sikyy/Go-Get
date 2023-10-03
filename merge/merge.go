package merge

import (
	"fmt"
	"io"
	"os"
)

// MergeChunks 合并下载的分片
func MergeChunks(chunkCount int64, mergedFileName string) {
	// 创建合并文件，存储最终合并文件的内容
	mergedFile, err := os.Create(mergedFileName)
	if err != nil {
		fmt.Println("Error creating merged file:", err)
		return
	}
	defer mergedFile.Close()

	// 遍历所有分片文件，逐一读取内容并写入合并文件
	for i := int64(0); i < chunkCount; i++ {
		chunkFileName := fmt.Sprintf("chunk_%d.tmp", i)
		chunkFile, err := os.Open(chunkFileName)
		if err != nil {
			fmt.Println("Error opening chunk file:", err)
			return
		}
		defer chunkFile.Close()

		_, err = io.Copy(mergedFile, chunkFile)
		if err != nil {
			fmt.Println("Error copying chunk to merged file:", err)
			return
		}

		// 删除已合并的分片文件
		os.Remove(chunkFileName)
	}
}
