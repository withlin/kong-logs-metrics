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
	Min [24]float64 `json:"min" binding:"required"`
	Max [24]float64 `json:"max" binding:"required"`
	Avg [24]float64 `json:"avg" binding:"required"`
}

// FindAggMetrics kong日志聚合统计Api  这是折线 条形 混住 图片
func FindAggMetrics(c *gin.Context) {

	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()
	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery()).Filter(elastic.NewRangeQuery("started_at").Gte("1524585600000").Lte("1524671999999").Format("epoch_millis"))
	// sou, err := query.Source()
	ctx := context.Background()
	avgAgg := elastic.NewAvgAggregation().Field("latencies.request")
	maxAgg := elastic.NewMaxAggregation().Field("latencies.request")
	minAgg := elastic.NewMinAggregation().Field("latencies.request")
	dataAgg := elastic.NewDateHistogramAggregation().Field("started_at").Interval("1h").TimeZone("Asia/Shanghai").MinDocCount(1).SubAggregation("avgAgg", avgAgg).SubAggregation("maxAgg", maxAgg).SubAggregation("minAgg", minAgg)

	searchResult, err := client.Search().Index("logstash-2018.04.25").Query(query).From(0).Size(0).Aggregation("DataAggs", dataAgg).Do(ctx)

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

// ConvertMap 赋值操作
func ConvertMap(arr []Bucket) (AggResult, error) {

	var min [24]float64
	var max [24]float64
	var avg [24]float64
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

	}
	result.Avg = avg
	result.Min = min
	result.Max = max
	return result, nil

}

func pieChar(c *gin.Context) {

	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery()).Filter(elastic.NewRangeQuery("started_at").Gte("1524585600000").Lte("1524671999999").Format("epoch_millis"))

	r1 := 30
	r2 := 90
	r3 := 120
	r4 := 150
	r5 := 180
	r6 := 210
	r7 := 240
	r8 := 270
	r9 := 300
	r10 := 500
	ctx := context.Background()
	rangeAgg := elastic.NewRangeAggregation().Field("latencies.request").AddRange(nil, r1).AddRange(r1, r2).AddRange(r2, r3).AddRange(r3, r4).AddRange(r4, r5).AddRange(r5, r6).AddRange(r6, r7).AddRange(r7, r8).AddRange(r8, r9).AddRange(r9, r10)
	tophitAgg := elastic.NewTopHitsAggregation().DocvalueFields("latencies.request").Size(2000).Sort("started_at", false)

	searchResult, err := client.Search().Index("logstash-2018.04.25").Query(query).Aggregation("rangeAgg", rangeAgg).Aggregation("tophitAgg", tophitAgg).Do(ctx)

	if err != nil {
		//do something
	}

	buf, err := json.Marshal(searchResult)
	if err != nil {
		//doSomthing
	}
	errCode := json.Unmarshal(buf, &aggMetrics)

	fmt.Println(errCode)
	c.JSON(http.StatusOK, gin.H{"message": "false", "data": "aaa"})

}


func  queryUrlName(c *gin.Context){
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery()).Filter(elastic.NewRangeQuery("started_at").Gte("1524585600000").Lte("1524671999999").Format("epoch_millis"))

}
