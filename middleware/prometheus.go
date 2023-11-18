package middleware

import (
	"Go-Get/metrics"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// 用于在接口被调用时递增相应的 QPS 计数器，并将 endpoint 作为 label。
func HandleEndpointQps() gin.HandlerFunc {
	return func(c *gin.Context) {
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
