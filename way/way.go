package way

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

//通用方法

// SendOutput 将输出发送到前端
func SendOutput(outputCh chan<- string, format string, args ...interface{}) {
	output := fmt.Sprintf(format, args...)
	fmt.Println(output) // 可选，用于在服务器端打印输出
	outputCh <- output
}

// 把json转换成bson
func JSONToBSONM(jsonData []byte) (bson.M, error) {
	var result bson.M
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}
	return result, nil
}
