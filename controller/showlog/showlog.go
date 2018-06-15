package showlog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"kong-logs-metrics/config"
	"kong-logs-metrics/controller/common"
	"kong-logs-metrics/model"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

//ShowLogs 展示日志
func ShowLogs(c *gin.Context) {
	userInter, _ := c.Get("user")
	user := userInter.(model.User)
	page := new(model.Page)
	SendErrJSON := common.SendErrJSON
	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery())

	ctx := context.Background()
	if err := c.ShouldBindJSON(&page); err == nil {
		if ok, _ := common.ES.IndexExists(page.DateValue).Do(ctx); ok {
			if page.PageNumber >= 0 && page.PageSize > 0 {
				page.Appid = user.AppID
				if user.Name == "admin" {

					AuthQuery(ctx, query, c, page)

				} else {
					phraseQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.headers.appid", page.Appid).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					AuthQuery(ctx, phraseQuery, c, page)

				}
			}
		} else {
			SendErrJSON("当前日期没有数据", c)
		}
	} else {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
}

//AuthQuery 根据不同的账号查询
func AuthQuery(ctx context.Context, query elastic.Query, c *gin.Context, page *model.Page) {
	SendErrJSON := common.SendErrJSON
	logs := new(model.Logs)
	searchResult, err := common.ES.Search().Index(page.DateValue).Type(config.Conf.ElasticSearch.LogStashType).Query(query).From(page.PageNumber).Size(page.PageSize).Do(ctx)
	if err != nil {
		SendErrJSON("ES查询错误", c)
		return
	}

	buf, err := json.Marshal(searchResult)
	if err != nil {
		SendErrJSON("error", c)
		return

	}
	errCode := json.Unmarshal(buf, &logs)
	if errCode != nil {

	}

	hits := logs.Hits
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  hits,
	})
}

//FindLogDetailByID 通过索引的ID 查找某个日志的详情
func FindLogDetailByID(c *gin.Context) {

	id := new(model.ID)
	SendErrJSON := common.SendErrJSON
	if err := c.ShouldBindJSON(&id); err == nil {

		if id.ID != "" && id.IndexName != "" {

			query := elastic.NewIdsQuery().Ids(id.ID)
			fmt.Println(query.Source())
			ctx := context.Background()
			searchResult, err := common.ES.Search().Index(id.IndexName).Type(config.Conf.ElasticSearch.LogStashType).Query(query).Do(ctx)
			if err != nil {
				SendErrJSON("error", c)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"errNo": model.ErrorCode.SUCCESS,
				"msg":   "success",
				"data":  searchResult.Hits,
			})

		} else {
			SendErrJSON("无效的ID或者索引名称", c)
		}
	} else {

		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}

}

//FindLogByAPINameAndDate FindLogByAPINameAndDate
func FindLogByAPINameAndDate(c *gin.Context) {

	api := new(model.API)
	logs := new(model.Logs)
	SendErrJSON := common.SendErrJSON
	ctx := context.Background()
	if err := c.ShouldBindJSON(&api); err == nil {

		if api.Name != "" && api.Data != "" {

			res, _ := common.ES.IndexExists(api.Data).Do(ctx)
			if res {
				var searchResult interface{}
				var err error
				if api.Appid != "" {
					boolQueryMatch := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.headers.appid", api.Appid).Slop(0).Boost(1), elastic.NewMatchPhraseQuery("request.uri", api.Name).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					boolQueryWrap := elastic.NewBoolQuery().Must(boolQueryMatch).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					filterQueryWrap := elastic.NewBoolQuery().Filter(boolQueryWrap).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					searchResult, err = common.ES.Search().Index(api.Data).Type(config.Conf.ElasticSearch.LogStashType).Query(filterQueryWrap).From(api.PageNumber).Size(api.PageSize).Do(ctx)
					fmt.Println(searchResult)
					fmt.Println(err)
				} else {
					boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.uri", api.Name).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					searchResult, err = common.ES.Search().Index(api.Data).Type(config.Conf.ElasticSearch.LogStashType).Query(macth).From(api.PageNumber).Size(api.PageSize).Do(ctx)
				}

				if err != nil {
					SendErrJSON("error", c)
					return
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {
					SendErrJSON("error", c)
					return

				}
				errCode := json.Unmarshal(buf, &logs)
				if errCode != nil {

				}
				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.SUCCESS,
					"msg":   "success",
					"data":  logs.Hits,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.ERROR,
					"msg":   "当前日期没有数据，请选择其他日期",
					"data":  logs.Hits,
				})
			}

		} else {
			res, _ := common.ES.IndexExists(api.Data).Do(ctx)
			if res {
				boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.headers.appid", api.Appid).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				searchResult, err := common.ES.Search().Index(api.Data).Type(config.Conf.ElasticSearch.LogStashType).Query(macth).From(api.PageNumber).Size(api.PageSize).Do(ctx)

				if err != nil {
					SendErrJSON("error", c)
					return
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {
					SendErrJSON("error", c)
					return

				}
				errCode := json.Unmarshal(buf, &logs)
				if errCode != nil {
					SendErrJSON("error", c)
				}
				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.SUCCESS,
					"msg":   "success",
					"data":  logs.Hits,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.ERROR,
					"msg":   "当前日期没有数据，请选择其他日期",
					"data":  logs.Hits,
				})
			}
		}
	} else {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
}
