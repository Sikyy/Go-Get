package getname

import (
	"errors"
	"net/url"
	"strings"
)

// func CleanFileName(fileName string, maxLength int) string {
// 	// 删除不合法的字符
// 	invalidChars := []rune{'\\', '/', ':', '*', '?', '"', '<', '>', '|'}
// 	fileName = strings.Map(func(r rune) rune {
// 		if unicode.IsControl(r) || unicode.IsSpace(r) {
// 			return '_'
// 		}
// 		for _, c := range invalidChars {
// 			if r == c {
// 				return '_'
// 			}
// 		}
// 		return r
// 	}, fileName)

// 	// 截取文件名，确保不超过最大长度
// 	if maxLength > 0 && len(fileName) > maxLength {
// 		fileName = fileName[:maxLength]
// 	}

// 	return fileName
// }

// 从内容配置中提取文件名
func ExtractFileNameFromContentDisposition(contentDisposition string) (string, error) {
	parts := strings.Split(contentDisposition, ";")
	for _, part := range parts {
		if strings.Contains(part, "filename=") {
			// 找到包含文件名的部分
			fileNamePart := strings.TrimSpace(part)
			fileName := strings.TrimPrefix(fileNamePart, "filename=")
			// 去除引号（如果有的话）
			fileName = strings.Trim(fileName, `"`)
			// 解码文件名
			decodedFileName, err := url.QueryUnescape(fileName)
			if err != nil {
				return "", err
			}
			return decodedFileName, nil
		}
	}
	return "", errors.New("File name not found in Content-Disposition header")
}
