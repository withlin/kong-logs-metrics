package main

import (
	"fmt"
	"kong-logs-metrics/config"
)

func main() {
	// config.InitJSON()
	config.InitAll()
	// config.InitElasticSearchConfig()
	fmt.Print(config.TestCinfig.URL)
}
