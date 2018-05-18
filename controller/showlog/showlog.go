package showlog

import (
	"context"
	"encoding/json"
	"fmt"
	"kong-logs-metrics/model"
	"net/http"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

//ShowLogs 展示日志
func ShowLogs(c *gin.Context) {
	page := new(model.Page)
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	query := elastic.NewBoolQuery().Must(elastic.NewMatchAllQuery())
	fmt.Println(query.Source())
	ctx := context.Background()
	if err := c.ShouldBindJSON(&page); err == nil {
		fmt.Println(page.PageNumber)
		fmt.Println(page.PageSize)
		fmt.Println(page.DateValue)
		if page.PageNumber > 0 && page.PageSize > 0 {

			searchResult, err := client.Search().Index(page.DateValue).Query(query).From(page.PageNumber).Size(page.PageSize).Do(ctx)

			if err != nil {
				//do Something

			}

			buf, err := json.Marshal(searchResult)
			if err != nil {
				//doSomthing
			}
			logs := new(model.Logs)
			errCode := json.Unmarshal(buf, &logs)

			if errCode != nil {
				//doSometing
			}

			c.JSON(http.StatusOK, gin.H{"message": "ok", "data": searchResult.Hits})

		} else {

			c.JSON(http.StatusOK, gin.H{"message": "false", "error": "PageSize和PageNumber必须大于零"})
		}
	} else {
		fmt.Println(err.Error())
	}
}

//FindLogDetailByID 通过索引的ID 查找某个日志的详情
func FindLogDetailByID(c *gin.Context) {

	id := new(model.ID)
	if err := c.ShouldBindJSON(&id); err == nil {

		if id.ID != "" && id.IndexName != "" {
			client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
			if err != nil {
				panic(err)
			}
			defer client.Stop()

			query := elastic.NewIdsQuery().Ids(id.ID)
			fmt.Println(query.Source())
			ctx := context.Background()
			searchResult, err := client.Search().Index(id.IndexName).Type("logs").Query(query).Do(ctx)
			fmt.Println(id.ID)

			buf, err := json.Marshal(searchResult)
			if err != nil {
				//doSomthing
			}
			logs := new(model.Logs)
			errCode := json.Unmarshal(buf, &logs)

			if errCode != nil {
				//doSometing
			}

			hits := searchResult.Hits

			c.JSON(http.StatusOK, gin.H{"message": "ok", "data": hits.Hits})

		} else {
			c.JSON(http.StatusOK, gin.H{"message": "false", "error": "无效的ID"})
		}
	} else {

		fmt.Println(err.Error())
	}

}

//FindLogByAPINameAndDate FindLogByAPINameAndDate
func FindLogByAPINameAndDate(c *gin.Context) {

	api := new(model.API)
	if err := c.ShouldBindJSON(&api); err == nil {
		fmt.Println(api.Name)
		fmt.Println(api.Data)

		if api.Name != "" && api.Data != "" {
			client, err := elastic.NewClient(elastic.SetURL("http://192.168.199.17:9200"), elastic.SetSniff(false))
			if err != nil {
				panic(err)
			}
			defer client.Stop()

			ctx := context.Background()

			res, _ := client.IndexExists(api.Data).Do(ctx)
			if res {
				boolQuery := elastic.NewBoolQuery().Must(elastic.NewMatchPhraseQuery("request.uri", api.Name).Slop(0).Boost(1)).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				macth := elastic.NewBoolQuery().Filter(boolQuery).DisableCoord(false).AdjustPureNegative(true).Boost(1)

				searchResult, err := client.Search().Index(api.Data).Type("logs").Query(macth).From(0).Size(200).Do(ctx)

				buf, err := json.Marshal(searchResult)
				if err != nil {
					//doSomthing
				}
				logs := new(model.Logs)
				errCode := json.Unmarshal(buf, &logs)

				if errCode != nil {
					//doSometing
				}

				// c.IndentedJSON()
				c.JSON(http.StatusOK, gin.H{"message": "ok", "data": logs.Hits})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "false", "data": "当前日期没有数据，请选择其他日期"})
			}

		} else {
			c.JSON(http.StatusOK, gin.H{"message": "false", "error": "..."})
		}
	} else {
		fmt.Println(err.Error())
	}
}
