package middleware

import (
	"Go-Get/metrics"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// 用于在接口被调用时递增相应的 QPS 计数器，并将 endpoint 作为 label。
func HandleEndpointQps() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取当前请求的接口路径
		endpoint := c.Request.URL.Path
		fmt.Println(endpoint)
		//递增 QPS 计数器，将当前请求的 endpoint 作为 label。然后，通过 Inc 方法递增 QPS 计数器。
		// 排除 /metrics 路由
		// if endpoint != "/metrics" {
		metrics.EndpointsQPSMonitor.With(prometheus.Labels{metrics.EndpointsDataSubsystem: endpoint}).Inc()
		// }
		c.Next()
	}
}

func HandleEndpointLantency() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求的接口路径
		endpoint := c.Request.URL.Path
		// 记录当前时间，作为请求处理开始的时间点
		start := time.Now()
		defer func(c *gin.Context) {
			// 计算请求处理耗时
			lantency := time.Since(start)
			// 将耗时直接转换为毫秒，保留三位小数
			lantencyFloat64 := float64(lantency.Nanoseconds()) / 1e6

			fmt.Printf("请求接口:%v，请求处理耗时：%f ms\n", endpoint, lantencyFloat64)
			// 将当前请求的 endpoint 作为 label。然后通过 Observe 方法记录耗时。
			metrics.EndpointsLantencyMonitor.With(prometheus.Labels{metrics.EndpointsDataSubsystem: endpoint}).Observe(lantencyFloat64)
		}(c)
		c.Next()
	}
}
