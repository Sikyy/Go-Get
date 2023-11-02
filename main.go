package main

import (
	"Go-Get/router"
)

func main() {
	// 初始化路由
	r := router.Router()
	r.Run(":9000")
}
