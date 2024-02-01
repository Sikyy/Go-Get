package metrics

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	PrometheusUrl          = "http://127.0.0.1:9091" //Prometheus PushGateway 的地址，用于推送指标数据。
	PrometheusJob          = "go_get_prometheus_qps" //Prometheus job 名称，用于标识这个应用程序的任务。,需要和yaml的job_name一致
	PrometheusNamespace    = "go_get_data"           //Prometheus 指标的命名空间。
	EndpointsDataSubsystem = "endpoints"             //Prometheus 指标的子系统，用于更细粒度地标识指标
)

var (
	Pusher *push.Pusher //Prometheus pusher，用于将指标数据推送到 Prometheus PushGateway。

	//一个 Counter 向量，用于记录接口 QPS 的统计信息.用于记录每个接口的请求数量。
	EndpointsQPSMonitor = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: PrometheusNamespace,    //命名空间
			Subsystem: EndpointsDataSubsystem, //子系统
			Name:      "QPS_statistic",        //指标名称
			Help:      "统计QPS数据",              //帮助信息
		}, []string{EndpointsDataSubsystem}, //用于标识不同的接口
	)
	//一个 Histogram 向量，用于记录接口耗时的统计信息。用于记录每个接口的耗时。
	EndpointsLantencyMonitor = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: PrometheusNamespace,                                      //命名空间
			Subsystem: EndpointsDataSubsystem,                                   //子系统
			Name:      "lantency_statistic",                                     //指标名称
			Help:      "统计耗时数据",                                                 //帮助信息
			Buckets:   []float64{1, 5, 10, 20, 50, 100, 500, 1000, 5000, 10000}, //耗时区间，如果一个请求的耗时在某个区间内，就会被记录到对应的桶中。
		}, []string{EndpointsDataSubsystem}, //用于标识不同的接口
	)
)

// 初始化指标
func Init() {
	// 创建一个新的pusher实例，使用PrometheusUrl和PrometheusJob进行配置
	Pusher = push.New(PrometheusUrl, PrometheusJob)

	// 使用prometheus.MustRegister函数注册两个监控指标，分别是EndpointsQPSMonitor和EndpointsLantencyMonitor
	prometheus.MustRegister(
		EndpointsQPSMonitor,
		EndpointsLantencyMonitor,
	)

	// 将EndpointsQPSMonitor和EndpointsLantencyMonitor添加到Pusher实例中的收集器列表中
	Pusher.Collector(EndpointsQPSMonitor)
	Pusher.Collector(EndpointsLantencyMonitor)
}

func PushGateway() {
	// 每15秒上报一次数据
	for range time.Tick(15 * time.Second) {
		if err := Pusher.Add(); err != nil {
			log.Println(err)
		}
		log.Println("push ")
	}
}
