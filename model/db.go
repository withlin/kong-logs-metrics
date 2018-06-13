package model

import (
	"fmt"
	"os"
	"time"

	"kong-logs-metrics/config"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

// DB  mysql数据库连接
var DB *gorm.DB

//RedisPool Redis连接池
var RedisPool *redis.Pool

func initDB() {
	// DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.Charset
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", config.Conf.Mysql.User, config.Conf.Mysql.Passworld, config.Conf.Mysql.Host, config.Conf.Mysql.Port, config.Conf.Mysql.Database, config.Conf.Mysql.Charset)

	db, err := gorm.Open(config.Conf.Mysql.Dialect, url)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if config.Conf.GoConf.Env == DevelopmentMode {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(config.Conf.Mysql.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.Conf.Mysql.MaxOpenConns)
	DB = db
}

func initRedis() {
	url := fmt.Sprintf("%s:%d", config.Conf.Redis.Host, config.Conf.Redis.Port)
	RedisPool = &redis.Pool{
		MaxIdle:     config.Conf.Redis.MaxIdle,
		MaxActive:   config.Conf.Redis.MaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url, redis.DialPassword(config.Conf.Redis.Password))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}
