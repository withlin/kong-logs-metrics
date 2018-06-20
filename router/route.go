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
		api.POST("/piechart", middleware.AuthUser, agg.PieChar)
		api.POST("/test/queryurlname", middleware.AuthUser, agg.QueryURLName)
		api.POST("/checklogin", login.PostCheckLogin)
		api.POST("/showlogs", middleware.AuthUser, showlog.ShowLogs)
		api.POST("/findlogdetailbyid", middleware.AuthUser, showlog.FindLogDetailByID)
		api.POST("/findlogsbyapiname", middleware.AuthUser, showlog.FindLogByAPINameAndDate)
		api.POST("findmatchid", middleware.AuthUser, agg.MatchID)
	}
}
