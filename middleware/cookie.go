// package middleware

// import (
// 	"fmt"
// 	"github.com/gin-gonic/gin"
// 	"kong-logs-metrics/config"
// )

// func RefreshTokenCookie(c *gin.Context) {
// 	tokenString, err:= c.Cookie("token")
// 	fmt.Println(err)
// 	if tokenString !="" && err == nil {
// 		c.SetCookie("token",tokenString,config.ServerConfig.TokenMaxAge,"/","",true,true)
// 		if user,err := getUser
// 	}
// }
