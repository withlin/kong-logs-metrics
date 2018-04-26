package main

import (
	"fmt"
	"io"
	"kong-logs-metrics/config"
	"os"

	"kong-logs-metrics/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// config.InitJSON()
	config.InitAll()
	// config.InitElasticSearchConfig()
	fmt.Print(config.TestCinfig.URL)
	fmt.Print(config.ServerConfig.APIPrefix)
	fmt.Print(config.ServerConfig.LogDir)
	fmt.Print(config.ServerConfig.LogFile)
	fmt.Print(config.ServerConfig.Port)
	gin.Default()

	fmt.Println("gin.Version: ", gin.Version)

	gin.SetMode(gin.DebugMode)

	logFile, err := os.OpenFile(config.ServerConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(-1)
	}
	gin.DefaultWriter = io.MultiWriter(logFile)

	app := gin.New()
	maxSize := int64(32)
	app.MaxMultipartMemory = maxSize << 20 // 3 MiB

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	app.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.Use(gin.Recovery())

	router.Route(app)
	// app.Run(":8888")
	app.Run(":" + fmt.Sprintf("%d", config.ServerConfig.Port))
}
