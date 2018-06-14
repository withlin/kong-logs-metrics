package login

import (
	"fmt"
	"net/http"

	"kong-logs-metrics/controller/common"

	"kong-logs-metrics/config"
	"kong-logs-metrics/model"

	"github.com/dgrijalva/jwt-go"
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
	var user model.User

	SendErrJSON := common.SendErrJSON
	if err := c.ShouldBindJSON(&loginCommand); err == nil {
		model.DB.Where("name = ? AND password = ?", loginCommand.Username, loginCommand.Password).First(&user)
		if user != (model.User{}) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id": user.ID,
			})
			tokenString, err := token.SignedString([]byte(config.Conf.GoConf.TokenSecret))
			fmt.Println(tokenString)
			if err != nil {
				fmt.Println(err.Error())
				SendErrJSON("内部错误", c)
				return
			}
			if err := model.UserToRedis(tokenString, user); err != nil {
				SendErrJSON("内部错误.", c)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"errNo": model.ErrorCode.SUCCESS,
				"msg":   "success",
				"data":  tokenString,
			})

		} else {
			SendErrJSON("账户或密码错误", c)
		}
	} else {
		SendErrJSON("error", c)
		return
	}
}
