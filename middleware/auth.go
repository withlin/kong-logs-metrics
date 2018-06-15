package middleware

import (
	"fmt"
	"kong-logs-metrics/controller/common"
	"kong-logs-metrics/model"

	"github.com/gin-gonic/gin"
)

//AuthUser 验证
func AuthUser(c *gin.Context) {
	var user model.User
	SendErrJSON := common.SendErrJSON
	tokenString := c.GetHeader("Access-Token")
	if tokenString == "" {
		SendErrJSON("未登录", c)
		return
	}

	var err error
	user, err = model.UserFromRedis(tokenString)
	fmt.Println(user)
	if err != nil {
		SendErrJSON("未登录", c)
		return
	}

	c.Set("user", user)
	c.Next()

}
