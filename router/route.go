package router

import (
	"kong-logs-metrics/config"
	"kong-logs-metrics/controller/elastic"
	"kong-logs-metrics/controller/login"
	"kong-logs-metrics/controller/showlog"

	"github.com/gin-gonic/gin"
)

// Route ?????api??
func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix)
	{
		api.GET("/test", agg.FindAggMetrics)
		api.GET("/test/PieChart", agg.PieChar)
		api.GET("/test/queryUrlName", agg.QueryURLName)
		api.POST("/checklogin", login.PostCheckLogin)
		api.GET("/showlogs", showlog.ShowLogs)
		api.POST("/findlogdetailbyid", showlog.FindLogDetailByID)
	}
}
