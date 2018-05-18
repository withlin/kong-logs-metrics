package agg

import (
	"context"
	"encoding/json"
	"fmt"
	"kong-logs-metrics/controller/common"
	"kong-logs-metrics/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

var aggMetrics model.AggMetrics

// FindAggMetrics kong日志聚合统计Api  这是折线 条形 混住 图片
func FindAggMetrics(c *gin.Context) {
	SendErrJSON := common.SendErrJSON
	loadaggchart := new(model.LoadAggChart)
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		SendErrJSON("error", c)
		return
	}
	defer client.Stop()

	if err := c.ShouldBindJSON(&loadaggchart); err == nil {
		if loadaggchart.LogstashName != "" {

			fmt.Println(loadaggchart.LogstashName)
			ctx := context.Background()

			if loadaggchart.Name != "" {
				boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.uri", loadaggchart.Name).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				avgAgg := elastic.NewAvgAggregation().Field("latencies.request")
				maxAgg := elastic.NewMaxAggregation().Field("latencies.request")
				minAgg := elastic.NewMinAggregation().Field("latencies.request")
				dataAgg := elastic.NewDateHistogramAggregation().Field("started_at").Interval("1h").TimeZone("Asia/Shanghai").MinDocCount(1).SubAggregation("avgAgg", avgAgg).SubAggregation("maxAgg", maxAgg).SubAggregation("minAgg", minAgg)

				searchResult, err := client.Search().Index(loadaggchart.LogstashName).Query(macth).From(0).Size(0).Aggregation("DataAggs", dataAgg).Do(ctx)

				if err != nil {
					//doSomething
					SendErrJSON("error", c)
					return
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {
					//doSomthing
					SendErrJSON("error", c)
					return

				}
				errCode := json.Unmarshal(buf, &aggMetrics)

				if errCode != nil {
					//doSometing
					SendErrJSON("error", c)
					return
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

				searchResult, err := client.Search().Index(loadaggchart.LogstashName).From(0).Size(0).Aggregation("DataAggs", dataAgg).Do(ctx)

				if err != nil {
					//doSomething
					SendErrJSON("error", c)
					return
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {
					//doSomthing
					SendErrJSON("error", c)
					return
				}
				errCode := json.Unmarshal(buf, &aggMetrics)

				if errCode != nil {
					//doSometing
					SendErrJSON("error", c)
					return
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
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()

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
				searchResult, err := client.Search().Index(piechartpost.LogstashName).Query(macth).Size(0).Aggregation("rangeAgg", rangeAgg).Do(ctx)

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
					SendErrJSON("error", c)
					return
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

				searchResult, err := client.Search().Index(piechartpost.LogstashName).Size(0).Aggregation("rangeAgg", rangeAgg).Do(ctx)

				if err != nil {
					SendErrJSON("error", c)
					return

					buf, err := json.Marshal(searchResult)
					if err != nil {
						//doSomthing
					}
					errCode := json.Unmarshal(buf, &pieMetrics)
					if errCode != nil {
						SendErrJSON("error", c)
						return
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
	url := new(model.URL)
	logstashname := new(model.DateValue)
	SendErrJSON := common.SendErrJSON
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		SendErrJSON("error", c)
		return
	}
	defer client.Stop()
	ctx := context.Background()
	if err := c.ShouldBindJSON(&logstashname); err == nil {
		if logstashname.LogstashName != "" {
			res, _ := client.IndexExists(logstashname.LogstashName).Do(ctx)

			if res {
				termAgg := elastic.NewTermsAggregation().Field("upstream_uri.keyword")

				searchResult, err := client.Search().Index(logstashname.LogstashName).Aggregation("termAgg", termAgg).Size(0).Do(ctx)

				if err != nil {
					SendErrJSON("error", c)
					return
				}

				da, err := json.Marshal(searchResult)
				if err != nil {
					SendErrJSON("error", c)
					return
				}

				err1 := json.Unmarshal(da, &url)
				if err1 != nil {
					//doSometing
					c.JSON(http.StatusOK, gin.H{"message": "false", "err": err1.Error()})
				}

				c.JSON(http.StatusOK, gin.H{"message": "ok", "data": url.Aggregations.TermAgg.Buckets})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "false", "data": "当前日期没有数据，请选择其他日期"})
			}

		}

	} else {
		SendErrJSON("error", c)
		return
	}
}
