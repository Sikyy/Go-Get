package way

import "fmt"

//通用方法

// SendOutput 将输出发送到前端
func SendOutput(outputCh chan<- string, format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	fmt.Println(output) // 可选，用于在服务器端打印输出
	outputCh <- output
}
