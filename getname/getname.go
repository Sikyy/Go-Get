package getname

import (
	"errors"
	"strings"
)

func ExtractFileNameFromContentDisposition(contentDisposition string) (string, error) {
	parts := strings.Split(contentDisposition, ";")
	for _, part := range parts {
		if strings.Contains(part, "filename=") {
			// 找到包含文件名的部分
			fileNamePart := strings.TrimSpace(part)
			fileName := strings.TrimPrefix(fileNamePart, "filename=")
			// 去除引号（如果有的话）
			fileName = strings.Trim(fileName, `"`)
			return fileName, nil
		}
	}
	return "", errors.New("File name not found in Content-Disposition header")
}
