package router

import (
	"kong-logs-metrics/controller/login"

	"kong-logs-metrics/config"
	"kong-logs-metrics/controller/elastic"
	"kong-logs-metrics/controller/showlog"
	"kong-logs-metrics/middleware"

	"github.com/gin-gonic/gin"
)

// Route ?????api??
func Route(router *gin.Engine) {
	apiPrefix := config.Conf.GoConf.APIPrefix

	api := router.Group(apiPrefix)
	{
		api.POST("/findaggmetrics", middleware.AuthUser, agg.FindAggMetrics)
		api.POST("/piechart", agg.PieChar)
		api.POST("/test/queryurlname", agg.QueryURLName)
		api.POST("/checklogin", login.PostCheckLogin)
		api.POST("/showlogs", middleware.AuthUser, showlog.ShowLogs)
		api.POST("/findlogdetailbyid", showlog.FindLogDetailByID)
		api.POST("/findlogsbyapiname", showlog.FindLogByAPINameAndDate)
		api.POST("findmatchid", agg.MatchID)
	}
}
