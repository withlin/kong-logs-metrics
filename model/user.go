package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"kong-logs-metrics/config"

	"github.com/gomodule/redigo/redis"
)

//User 用户表
type User struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	IsActive bool   `json:"isactive"`
	AppID    string `json:"appid"`
}

// UserFromRedis 从redis中取出用户信息
func UserFromRedis(tokenString string) (User, error) {

	RedisConn := RedisPool.Get()
	defer RedisConn.Close()

	userBytes, err := redis.Bytes(RedisConn.Do("GET", tokenString))
	if err != nil {
		fmt.Println(err)
		return User{}, errors.New("未登录")
	}
	var user User
	bytesErr := json.Unmarshal(userBytes, &user)
	if bytesErr != nil {
		fmt.Println(bytesErr)
		return user, errors.New("未登录")
	}
	return user, nil
}

// UserToRedis 将用户信息存到redis
func UserToRedis(tokenString string, user User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return errors.New("error")
	}

	RedisConn := RedisPool.Get()
	defer RedisConn.Close()

	if _, redisErr := RedisConn.Do("SET", tokenString, userBytes, "EX", config.Conf.GoConf.TokenMaxAge); redisErr != nil {
		fmt.Println("redis set failed: ", redisErr.Error())
		return errors.New("error")
	}
	return nil
}
