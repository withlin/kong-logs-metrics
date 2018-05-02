package test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

// AggMetrics 日志聚合
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
			Buckets []struct {
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
			} `json:"buckets"`
		} `json:"DataAggs"`
	} `json:"aggregations"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
}

var test AggMetrics

// AggSomething
func AggSomething(c *gin.Context) {
	// fmt.Println("\"message\":\"test\"")
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()
	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery()).Filter(elastic.NewRangeQuery("started_at").Gte("1524585600000").Lte("1524671999999").Format("epoch_millis"))
	// sou, err := query.Source()
	ctx := context.Background()
	avgAgg := elastic.NewAvgAggregation().Field("latencies.proxy")
	maxAgg := elastic.NewMaxAggregation().Field("latencies.proxy")
	minAgg := elastic.NewMinAggregation().Field("latencies.proxy")
	dataAgg := elastic.NewDateHistogramAggregation().Field("started_at").Interval("1h").TimeZone("Asia/Shanghai").MinDocCount(1).SubAggregation("avgAgg", avgAgg).SubAggregation("maxAgg", maxAgg).SubAggregation("minAgg", minAgg)

	searchResult, err := client.Search().Index("logstash-2018.04.25").Query(query).From(0).Size(0).Aggregation("DataAggs", dataAgg).Do(ctx)

	if err != nil {

	}
	if searchResult.TotalHits() != 1 {
		fmt.Errorf("expected Hits.TotalHits = %d; got: %d", 1, searchResult.TotalHits())
	}
	if _, found := searchResult.Aggregations["DataAggs"]; !found {
		fmt.Println("expected aggregation %q", "dhagg")
	}
	buf, err := json.Marshal(searchResult)
	a := json.Unmarshal(buf, &test)
	if a != nil {

	}
	fmt.Println("=========================%d\n============", test.Took)
	if err != nil {
		fmt.Println(err)
	}
	s := string(buf)
	fmt.Println(s)
	if i := strings.Index(s, `{"dhagg":{"buckets":[{"key_as_string":"2012-01-01`); i < 0 {
		fmt.Errorf("expected to serialize aggregation into string; got: %v", s)
	}
	c.JSON(200, searchResult)

}
