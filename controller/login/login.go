package login

import (
	"fmt"
	"net/http"

	"github.com/DevWithLin/kong-logs-metrics/controller/common"
	"github.com/gin-gonic/gin"
)

//User 模拟登录的对象
type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//PostCheckLogin 登录
func PostCheckLogin(c *gin.Context) {
	var loginCommand User
	fmt.Println("===========进来了============")
	SendErrJSON := common.SendErrJSON
	if err := c.ShouldBindJSON(&loginCommand); err == nil {

		if loginCommand.Username == "admin" && loginCommand.Password == "admin" {

			c.JSON(http.StatusOK, gin.H{"message": "ok", "data": "登录成功"})

		} else {
			c.JSON(http.StatusOK, gin.H{"message": "false", "data": "账户名无效或密码无效"})
		}
	} else {
		SendErrJSON("error", c)
		return
	}
}
