package test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

// latencies 延迟
type latencies struct {
	Request int `json:"request"`
	Proxy   int `json:"proxy"`
	Kong    int `json:"kong"`
}

// Hello 测试
func AggSomething(c *gin.Context) {
	// fmt.Println("\"message\":\"test\"")
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()
	query := elastic.NewBoolQuery().Must(elastic.1()).Filter(elastic.NewRangeQuery("started_at").Gte("1524585600000").Lte("1524671999999").Format("epoch_millis"))
	// sou, err := query.Source()
	ctx := context.Background()
	avgAgg := elastic.NewAvgAggregation().Field("latencies.proxy")
	dataAgg := elastic.NewDateHistogramAggregation().Field("started_at").Interval("1h").TimeZone("Asia/Shanghai").MinDocCount(1)
	maxAgg := elastic.NewMaxAggregation().Field("latencies.proxy")
	minAgg := elastic.NewMinAggregation().Field("latencies.proxy")
	test, err1 := client.Search().Index("logstash-2018.04.25").Query(query).From(0).Size(1).Aggregation("DataAggs", dataAgg).Aggregation("Avg-Proxy", avgAgg).Aggregation("Max-Agg", maxAgg).Aggregation("Min-Agg", minAgg).Do(ctx)

	xxx := test.Aggregations["Max-Agg"]
	var ar elastic.AggregationBucketKeyItems
	erraa := json.Unmarshal(*xxx, &ar)
	if erraa != nil {
		fmt.Printf("Unmarshal failed: %v\n", erraa)
		return
	}

	for _, item := range ar.Buckets {
		fmt.Printf("%v: %v\n", item.Key, item.DocCount)
	}


	for _, hit := range test.Hits.Hits {
		// hit.Index contains the name of the index

		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
		var lat latencies
		err := json.Unmarshal(*hit.Source, &lat)
		if err != nil {
			// Deserialization failed
		}
		c.JSON(200, "")
	}

	if test != nil {
		var lat latencies

		for _, item := range test.Each(reflect.TypeOf(lat)) {


			result, ok := item.(latencies)
			if ok {
				fmt.Printf("latencies by %s: %s\n", result.Proxy, result.Kong)
			} else {

				fmt.Errorf("错误啦啊啊啊 ")
			}
		}

		fmt.Println("tets========= %s \n=====", test.TotalHits())
	} else {
		fmt.Print(err1)
	}

	if err != nil {
		panic(err) 
	}

	
	data, err := json.Marshal(test)
	if err != nil {
		panic(err)
	}
	s := string(data)

	fmt.Println(s)

}
