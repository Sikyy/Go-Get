package getname

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// 正则表达式匹配动画名称，//一些很抽象的匹配不了
func ExtractAnimeName(filename string) string {
	// 使用正则表达式匹配动画名称
	re := regexp.MustCompile(`\[.*?]\s*([\p{L}\p{N}\s&-]+)\s*-`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) >= 2 {
		animeName := strings.TrimSpace(matches[1])
		return animeName
	}
	// 如果未匹配到，返回空字符串或默认值
	return ""
}

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
