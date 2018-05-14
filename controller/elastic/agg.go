package agg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

// AggMetrics 日志聚合对象
type AggMetrics struct {
	Took     int    `json:"took"`
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total    int           `json:"total"`
		MaxScore int           `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Suggest      interface{} `json:"suggest"`
	Aggregations struct {
		DataAggs struct {
			Buckets []Bucket `json:"buckets"`
		} `json:"DataAggs"`
	} `json:"aggregations"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

var aggMetrics AggMetrics

//Bucket Bucket 对象
type Bucket struct {
	KeyAsString time.Time `json:"key_as_string"`
	Key         int64     `json:"key"`
	DocCount    int       `json:"doc_count"`
	MaxAgg      struct {
		Value float64 `json:"value"`
	} `json:"maxAgg"`
	MinAgg struct {
		Value float64 `json:"value"`
	} `json:"minAgg"`
	AvgAgg struct {
		Value float64 `json:"value"`
	} `json:"avgAgg"`
}

//AggResult 聚合的结果
type AggResult struct {
	Min             [24]float64 `json:"min" binding:"required"`
	Max             [24]float64 `json:"max" binding:"required"`
	Avg             [24]float64 `json:"avg" binding:"required"`
	Count           [24]int     `json:"count" binding:"required"`
	TotalCount      int         `json:"totalCount" binding:"required"`
	ShareTotalCount int         `json:"shareTotalCount" binding:"required"`
}

//LoadAggChart post请求过来的数据
type LoadAggChart struct {
	LogstashName string `json:"logstastname" binding:"required`
	Name         string `json:"name" binding:"required`
}

// FindAggMetrics kong日志聚合统计Api  这是折线 条形 混住 图片
func FindAggMetrics(c *gin.Context) {
	loadaggchart := new(LoadAggChart)
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
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
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {
					//doSomthing
				}
				errCode := json.Unmarshal(buf, &aggMetrics)

				if errCode != nil {
					//doSometing
				}

				aggResult := aggMetrics.Aggregations.DataAggs.Buckets

				result, err := ConvertMap(aggResult)

				if err != nil {
					c.JSON(http.StatusOK, gin.H{"message": "false", "data": err})
				}

				// c.JSON(200, bbb)
				c.IndentedJSON(http.StatusOK, gin.H{"message": "ok", "data": result})
			} else {
				avgAgg := elastic.NewAvgAggregation().Field("latencies.request")
				maxAgg := elastic.NewMaxAggregation().Field("latencies.request")
				minAgg := elastic.NewMinAggregation().Field("latencies.request")
				dataAgg := elastic.NewDateHistogramAggregation().Field("started_at").Interval("1h").TimeZone("Asia/Shanghai").MinDocCount(1).SubAggregation("avgAgg", avgAgg).SubAggregation("maxAgg", maxAgg).SubAggregation("minAgg", minAgg)

				searchResult, err := client.Search().Index(loadaggchart.LogstashName).From(0).Size(0).Aggregation("DataAggs", dataAgg).Do(ctx)

				if err != nil {
					//doSomething
				}
				buf, err := json.Marshal(searchResult)
				if err != nil {
					//doSomthing
				}
				errCode := json.Unmarshal(buf, &aggMetrics)

				if errCode != nil {
					//doSometing
				}

				aggResult := aggMetrics.Aggregations.DataAggs.Buckets

				result, err := ConvertMap(aggResult)

				if err != nil {
					c.JSON(http.StatusOK, gin.H{"message": "false", "data": err})
				}

				// c.JSON(200, bbb)
				c.IndentedJSON(http.StatusOK, gin.H{"message": "ok", "data": result})
			}

		}
	}

}

// ConvertMap 赋值操作
func ConvertMap(arr []Bucket) (AggResult, error) {

	var min [24]float64
	var max [24]float64
	var avg [24]float64
	var count [24]int

	var result AggResult
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

//PieMetrics Pie聚合返回的实体
type PieMetrics struct {
	Took     int    `json:"took"`
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total    int           `json:"total"`
		MaxScore int           `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Suggest      interface{} `json:"suggest"`
	Aggregations struct {
		RangeAgg struct {
			Buckets []struct {
				Key      string  `json:"key"`
				To       float64 `json:"to"`
				DocCount int     `json:"doc_count"`
				From     float64 `json:"from,omitempty"`
			} `json:"buckets"`
		} `json:"rangeAgg"`
	} `json:"aggregations"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

//PieResult 返回Pie结果
type PieResult struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

var pieMetrics PieMetrics

// PieChar 圆表查询
func PieChar(c *gin.Context) {
	piechartpost := new(LoadAggChart)
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

				rangeAgg := elastic.NewRangeAggregation().Field("latencies.request").AddRange(nil, r1).AddRange(r1, r2).AddRange(r2, r3).AddRange(r3, r4).AddRange(r4, r5).AddRange(r5, r6).AddRange(r6, r7).AddRange(r7, r8).AddRange(r8, r9).AddRange(r9, r10).AddUnboundedFrom(r10)
				searchResult, err := client.Search().Index(piechartpost.LogstashName).Query(macth).Size(0).Aggregation("rangeAgg", rangeAgg).Do(ctx)

				if err != nil {
					//do something
					fmt.Println("发生错误了==================")
				}

				buf, err := json.Marshal(searchResult)
				if err != nil {
					//doSomthing
				}
				errCode := json.Unmarshal(buf, &pieMetrics)
				if errCode != nil {
					fmt.Println(errCode)
				}

				agg := pieMetrics.Aggregations
				rAgg := agg.RangeAgg
				pieBuckets := rAgg.Buckets
				ms := "ms"

				var item PieResult

				pieResults := []PieResult{}
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
					//do something
					fmt.Println("发生错误了==================" + err.Error())
				}

				buf, err := json.Marshal(searchResult)
				if err != nil {
					//doSomthing
				}
				errCode := json.Unmarshal(buf, &pieMetrics)
				if errCode != nil {
					fmt.Println(errCode)
				}

				agg := pieMetrics.Aggregations
				rAgg := agg.RangeAgg
				pieBuckets := rAgg.Buckets
				ms := "ms"

				var item PieResult

				pieResults := []PieResult{}
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
			}

		}
	}

}

//URL 查询urlname
type URL struct {
	Took     int    `json:"took"`
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total    int           `json:"total"`
		MaxScore int           `json:"max_score"`
		Hits     []interface{} `json:"hits"`
	} `json:"hits"`
	Suggest      interface{} `json:"suggest"`
	Aggregations struct {
		TermAgg struct {
			DocCountErrorUpperBound int             `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int             `json:"sum_other_doc_count"`
			Buckets                 []ResultBuckets `json:"buckets"`
		} `json:"termAgg"`
	} `json:"aggregations"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

//ResultBuckets bucketResult
type ResultBuckets struct {
	Key string `json:"key"`
	// DocCount int    `json:"doc_count"`
}

type DateValue struct {
	LogstashName string `json:"logstastname" binding:"required`
}

//QueryURLName 查询请求的API名称
func QueryURLName(c *gin.Context) {
	url := new(URL)
	logstashname := new(DateValue)
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
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
					//doSomething
					fmt.Println(err.Error())
				}

				da, err := json.Marshal(searchResult)
				if err != nil {
					//doSometing
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
		fmt.Println("发生错误啦==================")
	}

}
