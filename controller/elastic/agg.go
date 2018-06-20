package agg

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"kong-logs-metrics/config"
	"kong-logs-metrics/controller/common"
	"kong-logs-metrics/model"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

var aggMetrics model.AggMetrics

// FindAggMetrics kong日志聚合统计Api  这是折线 条形 混住 图片
func FindAggMetrics(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	loadaggchart := new(model.LoadAggChart)
	if err := c.ShouldBindJSON(&loadaggchart); err == nil {
		if loadaggchart.LogstashName != "" {
			ctx := context.Background()

			if loadaggchart.Name != "" {
				boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.uri", loadaggchart.Name).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				avgAgg := elastic.NewAvgAggregation().Field("latencies.request")
				maxAgg := elastic.NewMaxAggregation().Field("latencies.request")
				minAgg := elastic.NewMinAggregation().Field("latencies.request")
				dataAgg := elastic.NewDateHistogramAggregation().Field("started_at").Interval("1h").TimeZone("Asia/Shanghai").MinDocCount(1).SubAggregation("avgAgg", avgAgg).SubAggregation("maxAgg", maxAgg).SubAggregation("minAgg", minAgg)

				searchResult, err := common.ES.Search().Index(loadaggchart.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Query(macth).From(0).Size(0).Aggregation("DataAggs", dataAgg).Do(ctx)

				if err != nil {

					SendErrJSON("error", c)
					return
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {

					SendErrJSON("error", c)
					return

				}
				errCode := json.Unmarshal(buf, &aggMetrics)

				if errCode != nil {

					// SendErrJSON("error", c)
					// return
				}

				aggResult := aggMetrics.Aggregations.DataAggs.Buckets

				result, err := ConvertMap(aggResult)

				if err != nil {
					c.JSON(http.StatusOK, gin.H{"message": "false", "data": err})
				}

				c.IndentedJSON(http.StatusOK, gin.H{"message": "ok", "data": result})
			} else {
				avgAgg := elastic.NewAvgAggregation().Field("latencies.request")
				maxAgg := elastic.NewMaxAggregation().Field("latencies.request")
				minAgg := elastic.NewMinAggregation().Field("latencies.request")
				dataAgg := elastic.NewDateHistogramAggregation().Field("started_at").Interval("1h").TimeZone("Asia/Shanghai").MinDocCount(1).SubAggregation("avgAgg", avgAgg).SubAggregation("maxAgg", maxAgg).SubAggregation("minAgg", minAgg)

				searchResult, err := common.ES.Search().Index(loadaggchart.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).From(0).Size(0).Aggregation("DataAggs", dataAgg).Do(ctx)

				if err != nil {

					SendErrJSON("error", c)
					return
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {

					SendErrJSON("error", c)
					return
				}
				errCode := json.Unmarshal(buf, &aggMetrics)

				if errCode != nil {

					// SendErrJSON("error", c)
					// return
				}

				aggResult := aggMetrics.Aggregations.DataAggs.Buckets

				result, err := ConvertMap(aggResult)

				if err != nil {
					c.JSON(http.StatusOK, gin.H{"message": "false", "data": err})
				}

				c.IndentedJSON(http.StatusOK, gin.H{"message": "ok", "data": result})
			}

		}
	}

}

// ConvertMap 赋值操作
func ConvertMap(arr []model.Bucket) (model.AggResult, error) {

	var min [24]float64
	var max [24]float64
	var avg [24]float64
	var count [24]int

	var result model.AggResult
	for _, elem := range arr {
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return result, err
		}
		localResult := elem.KeyAsString.In(location)
		tim, err := strconv.Atoi(localResult.Format("15"))
		if err != nil {
			return result, err
		}
		min[tim] = elem.MinAgg.Value
		max[tim] = elem.MaxAgg.Value
		avg[tim] = elem.AvgAgg.Value
		count[tim] = elem.DocCount

	}
	result.Avg = avg
	result.Min = min
	result.Max = max
	result.Count = count
	result.TotalCount = aggMetrics.Hits.Total
	result.ShareTotalCount = aggMetrics.Shards.Total
	return result, nil

}

var pieMetrics model.PieMetrics

// PieChar 圆表查询
func PieChar(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	piechartpost := new(model.LoadAggChart)
	r1 := 2000
	r2 := 4000
	r3 := 6000
	r4 := 8000
	r5 := 10000
	r6 := 12000
	r7 := 14000
	r8 := 16000
	r9 := 18000
	r10 := 20000
	ctx := context.Background()
	if err := c.ShouldBindJSON(&piechartpost); err == nil {

		if piechartpost.LogstashName != "" {
			if piechartpost.Name != "" {
				boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.uri", piechartpost.Name).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				rangeAgg := elastic.NewRangeAggregation().Field("latencies.request").AddRange(nil, r1).AddRange(r1, r2).AddRange(r2, r3).AddRange(r3, r4).AddRange(r4, r5).AddRange(r5, r6).AddRange(r6, r7).AddRange(r7, r8).AddRange(r8, r9).AddRange(r9, r10).AddUnboundedFrom(10000000000)
				searchResult, err := common.ES.Search().Index(piechartpost.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Query(macth).Size(0).Aggregation("rangeAgg", rangeAgg).Do(ctx)

				if err != nil {
					SendErrJSON("error", c)
					return
				}

				buf, err := json.Marshal(searchResult)
				if err != nil {
					SendErrJSON("error", c)
					return

				}
				errCode := json.Unmarshal(buf, &pieMetrics)
				if errCode != nil {

				}

				agg := pieMetrics.Aggregations
				rAgg := agg.RangeAgg
				pieBuckets := rAgg.Buckets
				ms := "ms"

				var item model.PieResult

				pieResults := []model.PieResult{}
				// fmt.Println(PieBuckets[1].DocCount)
				for _, elem := range pieBuckets {

					if elem.From == 0 {
						// fmt.Println(index)
						to := strconv.FormatFloat(elem.To, 'f', 0, 64)
						item.Name = to + ms
						item.Value = elem.DocCount
						pieResults = append(pieResults, item)

					} else {
						// fmt.Println(index)
						to := strconv.FormatFloat(elem.To, 'f', 0, 64)
						from := strconv.FormatFloat(elem.From, 'f', 0, 64)
						item.Name = from + ms + "-" + to + ms
						item.Value = elem.DocCount
						pieResults = append(pieResults, item)

					}

				}

				c.JSON(http.StatusOK, gin.H{"message": "ok", "data": pieResults})
				// c.JSON(http.StatusOK, searchResult)

			} else {
				rangeAgg := elastic.NewRangeAggregation().Field("latencies.request").AddRange(nil, r1).AddRange(r1, r2).AddRange(r2, r3).AddRange(r3, r4).AddRange(r4, r5).AddRange(r5, r6).AddRange(r6, r7).AddRange(r7, r8).AddRange(r8, r9).AddRange(r9, r10).AddRange(r10, nil)
				// tophitAgg := elastic.NewTopHitsAggregation().DocvalueFields("latencies.request").Sort("started_at", false)

				searchResult, err := common.ES.Search().Index(piechartpost.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Size(0).Aggregation("rangeAgg", rangeAgg).Do(ctx)

				if err == nil {

					buf, err := json.Marshal(searchResult)
					if err != nil {
						SendErrJSON("error", c)
						return
					}
					errCode := json.Unmarshal(buf, &pieMetrics)
					if errCode != nil {

					}

					agg := pieMetrics.Aggregations
					rAgg := agg.RangeAgg
					pieBuckets := rAgg.Buckets
					ms := "ms"

					var item model.PieResult

					pieResults := []model.PieResult{}

					for _, elem := range pieBuckets {

						if elem.From == 0 {

							to := strconv.FormatFloat(elem.To, 'f', 0, 64)
							item.Name = to + ms
							item.Value = elem.DocCount
							pieResults = append(pieResults, item)

						} else {

							to := strconv.FormatFloat(elem.To, 'f', 0, 64)
							from := strconv.FormatFloat(elem.From, 'f', 0, 64)
							item.Name = from + ms + "-" + to + ms
							item.Value = elem.DocCount
							pieResults = append(pieResults, item)

						}

					}

					c.JSON(http.StatusOK, gin.H{"message": "ok", "data": pieResults})

				}

			}
		}

	}
}

//QueryURLName 查询请求的API名称
func QueryURLName(c *gin.Context) {
	userInter, _ := c.Get("user")
	user := userInter.(model.User)
	// wapperQuery("upstream_uri.keyword", c, user)
	url := new(model.URL)

	logstashname := new(model.DateValue)
	SendErrJSON := common.SendErrJSON
	ctx := context.Background()

	if err := c.ShouldBindJSON(&logstashname); err == nil {

		if logstashname.LogstashName != "" {
			res, _ := common.ES.IndexExists(logstashname.LogstashName).Do(ctx)

			if res {
				var result []byte
				if user.Name == "admin" {
					var err error
					termAgg := elastic.NewTermsAggregation().Field("upstream_uri.keyword")
					searchResult, err := common.ES.Search().Index(logstashname.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Aggregation("termAgg", termAgg).Size(0).Do(ctx)

					if err != nil {
						SendErrJSON("error", c)
						return
					}

					result, err = json.Marshal(searchResult)
					if err != nil {
						SendErrJSON("error", c)
						return
					}
				} else {
					boolQueryMatch := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.headers.appid", user.AppID).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					boolQueryWrap := elastic.NewBoolQuery().Must(boolQueryMatch).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					filterQueryWrap := elastic.NewBoolQuery().Filter(boolQueryWrap).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					termAgg := elastic.NewTermsAggregation().Field("upstream_uri.keyword")
					searchResult, err := common.ES.Search().Index(logstashname.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Query(filterQueryWrap).Aggregation("termAgg", termAgg).Size(0).Do(ctx)

					if err != nil {
						SendErrJSON("error", c)
						return
					}

					result, err = json.Marshal(searchResult)
					if err != nil {
						SendErrJSON("error", c)
						return
					}
				}

				errCode := json.Unmarshal(result, &url)
				if errCode != nil {
					SendErrJSON("error", c)
				}
				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.SUCCESS,
					"msg":   "success",
					"data":  url.Aggregations.TermAgg.Buckets,
				})
			} else {

				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.ERROR,
					"msg":   "error",
					"data":  "当前日期没有数据，请选择其他日期",
				})
			}

		}

	} else {
		SendErrJSON("error", c)
		return
	}
}

func wapperQuery(keyword string, c *gin.Context, user model.User) {
	url := new(model.URL)

	logstashname := new(model.DateValue)
	SendErrJSON := common.SendErrJSON
	ctx := context.Background()

	if err := c.ShouldBindJSON(&logstashname); err == nil {

		if logstashname.LogstashName != "" {
			res, _ := common.ES.IndexExists(logstashname.LogstashName).Do(ctx)

			if res {
				var result []byte
				if user.Name == "admin" {
					var err error
					termAgg := elastic.NewTermsAggregation().Field(keyword)
					searchResult, err := common.ES.Search().Index(logstashname.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Aggregation("termAgg", termAgg).Size(0).Do(ctx)

					if err != nil {
						SendErrJSON("error", c)
						return
					}

					result, err = json.Marshal(searchResult)
					if err != nil {
						SendErrJSON("error", c)
						return
					}
				} else {
					boolQueryMatch := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery(keyword, user.AppID).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					boolQueryWrap := elastic.NewBoolQuery().Must(boolQueryMatch).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					filterQueryWrap := elastic.NewBoolQuery().Filter(boolQueryWrap).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					termAgg := elastic.NewTermsAggregation().Field(keyword)
					searchResult, err := common.ES.Search().Index(logstashname.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Query(filterQueryWrap).Aggregation("termAgg", termAgg).Size(0).Do(ctx)

					if err != nil {
						SendErrJSON("error", c)
						return
					}

					result, err = json.Marshal(searchResult)
					if err != nil {
						SendErrJSON("error", c)
						return
					}
				}

				errCode := json.Unmarshal(result, &url)
				if errCode != nil {
					SendErrJSON("error", c)
				}
				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.SUCCESS,
					"msg":   "success",
					"data":  url.Aggregations.TermAgg.Buckets,
				})
			} else {

				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.ERROR,
					"msg":   "error",
					"data":  "当前日期没有数据，请选择其他日期",
				})
			}

		}

	} else {
		SendErrJSON("error", c)
		return
	}
}

//MatchID 查询matchid并且去重
func MatchID(c *gin.Context) {
	userInter, _ := c.Get("user")
	user := userInter.(model.User)
	// wapperQuery("request.headers.appid.keyword", c, user)
	url := new(model.URL)
	keyword := "request.headers.appid.keyword"
	logstashname := new(model.DateValue)
	SendErrJSON := common.SendErrJSON
	ctx := context.Background()

	if err := c.ShouldBindJSON(&logstashname); err == nil {

		if logstashname.LogstashName != "" {
			res, _ := common.ES.IndexExists(logstashname.LogstashName).Do(ctx)

			if res {
				var result []byte
				if user.Name == "admin" {
					var err error
					termAgg := elastic.NewTermsAggregation().Field(keyword)
					searchResult, err := common.ES.Search().Index(logstashname.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Aggregation("termAgg", termAgg).Size(0).Do(ctx)

					if err != nil {
						SendErrJSON("error", c)
						return
					}

					result, err = json.Marshal(searchResult)
					if err != nil {
						SendErrJSON("error", c)
						return
					}
				} else {
					boolQueryMatch := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.headers.appid", user.AppID).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					boolQueryWrap := elastic.NewBoolQuery().Must(boolQueryMatch).DisableCoord(false).AdjustPureNegative(true).Boost(1)
					filterQueryWrap := elastic.NewBoolQuery().Filter(boolQueryWrap).DisableCoord(false).AdjustPureNegative(true).Boost(1)

					termAgg := elastic.NewTermsAggregation().Field(keyword)
					searchResult, err := common.ES.Search().Index(logstashname.LogstashName).Type(config.Conf.ElasticSearch.LogStashType).Query(filterQueryWrap).Aggregation("termAgg", termAgg).Size(0).Do(ctx)

					if err != nil {
						SendErrJSON("error", c)
						return
					}

					result, err = json.Marshal(searchResult)
					if err != nil {
						SendErrJSON("error", c)
						return
					}
				}

				errCode := json.Unmarshal(result, &url)
				if errCode != nil {
					SendErrJSON("error", c)
				}
				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.SUCCESS,
					"msg":   "success",
					"data":  url.Aggregations.TermAgg.Buckets,
				})
			} else {

				c.JSON(http.StatusOK, gin.H{
					"errNo": model.ErrorCode.ERROR,
					"msg":   "error",
					"data":  "当前日期没有数据，请选择其他日期",
				})
			}

		}

	} else {
		SendErrJSON("error", c)
		return
	}
}
