package common

import (
	"fmt"
	"kong-logs-metrics/config"
	"os"
	"sync"

	elastic "gopkg.in/olivere/elastic.v5"
)

var ES *elastic.Client
var once sync.Once

//ESClient es连接
func ESClient() *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(config.ESCinfig.Host), elastic.SetSniff(config.ESCinfig.SetSniff))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	once.Do(func() {
		ES = client
	})
	return ES
}

func init() {
	ESClient()
}
