package middleware

import (
	"fmt"
	"kong-logs-metrics/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AuthUser 验证
func AuthUser(c *gin.Context) {
	var user model.User
	tokenString := c.GetHeader("Access-Token")
	if tokenString == "" {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.LoginError,
			"msg":   "success",
			"data":  "",
		})
		return
	}

	var err error
	user, err = model.UserFromRedis(tokenString)
	fmt.Println(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.LoginError,
			"msg":   "success",
			"data":  "",
		})
		return
	}

	c.Set("user", user)
	c.Next()

}
