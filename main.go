package main

import (
	"Go-Get/metrics"
	"Go-Get/router"
)

func main() {

	// 初始化指标
	metrics.Init()

	// 初始化路由
	r := router.Router()

	// 每15秒上报一次数据给 PushGateway
	go metrics.PushGateway()

	r.Run(":9000")
}
