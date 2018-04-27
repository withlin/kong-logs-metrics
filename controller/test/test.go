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
func Hello(c *gin.Context) {
	// fmt.Println("\"message\":\"test\"")
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()
	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery()).Filter(elastic.NewRangeQuery("started_at").Gte("1524585600000").Lte("1524671999999").Format("epoch_millis"))
	sou, err := query.Source()
	ctx := context.Background()
	test, err1 := client.Search().Index("logstash-2018.04.25").Query(query).From(0).Size(1).Pretty(true).Do(ctx)

	for _, hit := range test.Hits.Hits {
		// hit.Index contains the name of the index

		// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
		var lat latencies
		err := json.Unmarshal(*hit.Source, &lat)
		if err != nil {
			// Deserialization failed
		}

		// Work with tweet
		// fmt.Print(string(*hit.Source))
		str := string(*hit.Source)
		c.JSON(200, str)
	}

	if test != nil {
		var lat latencies
		// for _;item :=  test.Each(reflect.TypeOf(lat)) {
		// 	if t, ok := item.(latencies); ok {
		// 		fmt.Printf("latencies by %s: %s\n", t.Request, t.Proxy)
		// 	}
		// }
		for _, item := range test.Each(reflect.TypeOf(lat)) {
			// if t, ok := item.(latencies); ok {
			// 	fmt.Printf("latencies by %s: %s\n", t.Proxy, t.Kong)
			// }
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

	// if query != nil {
	// 	fmt.Println(query)
	// }
	data, err := json.Marshal(sou)
	if err != nil {
		panic(err)
	}
	s := string(data)

	fmt.Println(s)

}
