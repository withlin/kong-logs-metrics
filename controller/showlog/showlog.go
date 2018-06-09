package showlog

import (
	"context"
	"encoding/json"
	"fmt"
	"kong-logs-metrics/config"
	"kong-logs-metrics/controller/common"
	"kong-logs-metrics/model"
	"net/http"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

//ShowLogs 展示日志
func ShowLogs(c *gin.Context) {
	page := new(model.Page)
	SendErrJSON := common.SendErrJSON
	logs := new(model.Logs)
	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery())
	ctx := context.Background()
	if err := c.ShouldBindJSON(&page); err == nil {
		if page.PageNumber > 0 && page.PageSize > 0 {
			fmt.Println(config.ESCinfig.LogstashType)
			searchResult, err := common.ES.Search().Index(page.DateValue).Type(config.ESCinfig.LogstashType).Query(query).From(page.PageNumber).Size(page.PageSize).Do(ctx)

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

			c.JSON(http.StatusOK, gin.H{"message": "ok", "data": logs.Hits})

		} else {

			c.JSON(http.StatusOK, gin.H{"message": "false", "error": "PageSize和PageNumber必须大于零"})
		}
	} else {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
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
			searchResult, err := common.ES.Search().Index(id.IndexName).Type(config.ESCinfig.LogstashType).Query(query).Do(ctx)
			if err != nil {
				SendErrJSON("error", c)
				return
			}

			hits := searchResult.Hits

			c.JSON(http.StatusOK, gin.H{"message": "ok", "data": hits.Hits})

		} else {
			c.JSON(http.StatusOK, gin.H{"message": "false", "error": "无效的ID"})
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

					searchResult, err = common.ES.Search().Index(api.Data).Type(config.ESCinfig.LogstashType).Query(filterQueryWrap).From(api.PageNumber).Size(api.PageSize).Do(ctx)
					fmt.Println(searchResult)
					fmt.Println(err)
				} else {
					boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.uri", api.Name).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					searchResult, err = common.ES.Search().Index(api.Data).Type(config.ESCinfig.LogstashType).Query(macth).From(api.PageNumber).Size(api.PageSize).Do(ctx)
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
				c.JSON(http.StatusOK, gin.H{"message": "ok", "data": logs.Hits})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "false", "data": "当前日期没有数据，请选择其他日期"})
			}

		} else {
			fmt.Println("==============只有AppId的查询进来了==========")
			res, _ := common.ES.IndexExists(api.Data).Do(ctx)
			if res {
				boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.headers.appid", api.Appid).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				searchResult, err := common.ES.Search().Index(api.Data).Type(config.ESCinfig.LogstashType).Query(macth).From(api.PageNumber).Size(api.PageSize).Do(ctx)

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
				c.JSON(http.StatusOK, gin.H{"message": "ok", "data": logs.Hits})
			}
		}
	} else {
		fmt.Println(err.Error())
		SendErrJSON("error", c)
		return
	}
}

// //FindLogByAppid 通过Appid去匹配日志
// func FindLogByAppid(c *gin.Context) {
// 	matchid := new(model.MatchAppID)

// 	SendErrJSON := common.SendErrJSON
// 	if err := c.ShouldBindJSON(&matchid); err == nil {
// 		if matchid.Appid != "" {
// 			ctx := context.Background()

// 			res, _ := common.ES.IndexExists(matchid.Data).Do(ctx)

// 			if res {
// 				boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.headers.appid", matchid.Appid).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)
// 				filterQuery := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)
// 				searchResult, err := common.ES.Search().Index(matchid.Data).Type(config.ESCinfig.LogstashType).Query(filterQuery).From(matchid.PageNumber).Size(matchid.PageSize).Do(ctx)

// 				if err != nil {
// 					SendErrJSON("error", c)
// 					return
// 				}
// 				buf, err := json.Marshal(searchResult)
// 				if err != nil {
// 					SendErrJSON("error", c)
// 					return

// 				}
// 				errCode := json.Unmarshal(buf, &logs)
// 				if errCode != nil {

// 				}
// 				c.JSON(http.StatusOK, gin.H{"message": "ok", "data": logs.Hits})
// 			}
// 		}
// 	}
// }
