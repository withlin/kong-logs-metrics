package login

import (
	"fmt"
	"net/http"

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
	if err := c.ShouldBindJSON(&loginCommand); err == nil {
		fmt.Println("========================" + loginCommand.Password)
		if loginCommand.Username == "" || loginCommand.Password == "" {
			c.JSON(http.StatusOK, gin.H{"message": "false", "data": "无效账户名或密码"})
		} else if loginCommand.Username == "admin" && loginCommand.Password == "admin" {
			c.JSON(http.StatusOK, gin.H{"message": "ok", "data": "登录成功"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "false", "data": "账户名无效或密码无效"})
		}
	} else {
		fmt.Println(err.Error())
	}
}
