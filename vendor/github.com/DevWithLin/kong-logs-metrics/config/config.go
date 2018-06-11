package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"unicode/utf8"

	"github.com/DevWithLin/kong-logs-metrics/utils"

	yaml "gopkg.in/yaml.v2"
)

var jsonData map[string]interface{}

// ESConfi
type conf struct {
	ElasticSearch struct {
		Host         string `yaml:"host"`
		SetSniff     bool   `yaml:"setsniff"`
		LogStashType string `yaml:"logstashtype"`
		URL          string
	} `yaml:"elasticsearch"`
	GoConf struct {
		APIPrefix   string `yaml:"host"`
		Port        int    `yaml:"port"`
		TokenMaxAge int    `yaml:"tokenmaxage"`
		Env         string `yaml:"env"`
		LogDir      string `yaml:"logdir"`
	} `yaml:"go"`
}

//Conf 配置文件
var Conf conf

func initYaml() {
	bytes, err := ioutil.ReadFile("./config.yml")
	fmt.Println(bytes)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		os.Exit(-1)
	}
	err = yaml.Unmarshal(bytes, &Conf)
	if err != nil {
		fmt.Println("Yml File Unmarshal :", err.Error())
	}
}

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
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	URL          string `yaml:"URL"`
	SetSniff     bool   `yaml:"setsniff"`
	LogstashType string `yaml:"logstashtype"`
}

// ESCinfig 相关测试配置
var ESCinfig elasticSearchConfig

//InitElasticSearchConfig  相关配置
func initElasticSearchConfig() {
	utils.SetStructByJSON(&ESCinfig, jsonData["elasticsearch"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%d", ESCinfig.Host, ESCinfig.Port)
	ESCinfig.URL = url
}

type serverConfig struct {
	LogDir      string
	APIPrefix   string
	Port        int
	TokenMaxAge int
	LogFile     string
	Env         string
}

// ServerConfig 服务端配置
var ServerConfig serverConfig

func initServerConfig() {
	utils.SetStructByJSON(&ServerConfig, jsonData["go"].(map[string]interface{}))
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
func init() {
	initYaml()
}
