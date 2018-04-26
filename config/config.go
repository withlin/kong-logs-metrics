package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kong-logs-metrics/utils"
	"os"
	"regexp"
	"unicode/utf8"
)

var jsonData map[string]interface{}

//initJSON 初始化相关config.json相关数据
func initJSON() {
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		os.Exit(-1)
	}

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)

	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		fmt.Println("invalid config: ", err.Error())
		os.Exit(-1)
	}
}

type elasticSearchConfig struct {
	Host string
	Port int
	URL  string
}

// TestCinfig 相关测试配置
var TestCinfig elasticSearchConfig

//InitElasticSearchConfig  相关配置
func initElasticSearchConfig() {
	utils.SetStructByJSON(&TestCinfig, jsonData["elasticsearch"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%d", TestCinfig.Host, TestCinfig.Port)
	TestCinfig.URL = url
}

type serverConfig struct {
	LogDir      string
	APIPrefix   string
	Port        int
	TokenMaxAge int8
	LogFile     string
}

// ServerConfig 服务端配置
var ServerConfig serverConfig

func initServerConfig() {
	utils.SetStructByJSON(&TestCinfig, jsonData["go"].(map[string]interface{}))
	sep := string(os.PathSeparator)
	execPath, _ := os.Getwd()
	length := utf8.RuneCountInString(execPath)
	lastChar := execPath[length-1:]
	if lastChar != sep {
		execPath = execPath + sep
	}

	ymdStr := utils.GetTodayYMD("-")

	if ServerConfig.LogDir == "" {
		ServerConfig.LogDir = execPath
	} else {
		length := utf8.RuneCountInString(ServerConfig.LogDir)
		lastChar := ServerConfig.LogDir[length-1:]
		if lastChar != sep {
			ServerConfig.LogDir = ServerConfig.LogDir + sep
		}
	}
	ServerConfig.LogFile = ServerConfig.LogDir + ymdStr + ".log"
}

//InitAll 初始化全部的数据
func InitAll() {
	initJSON()
	initElasticSearchConfig()
	initServerConfig()
}
