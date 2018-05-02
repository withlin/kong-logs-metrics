package router

import (
	"kong-logs-metrics/config"
	"kong-logs-metrics/controller/elastic"

	"github.com/gin-gonic/gin"
)

// Route 路由
func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix)
	{
		api.GET("/test", test.AggMetricsController)
	}
}
