package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

// Configure prometheus metrics
// Enables gin metrics on the gin server
func Init(router *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	// set duration buckets
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)
}
