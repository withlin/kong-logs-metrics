package middleware

// import (
// 	"fmt"

// 	"kong-logs-metrics/model"

// 	"github.com/gin-gonic/gin"
// )

// //RefreshTokenCookie 每次请求进来就刷新Token
// func RefreshTokenCookie(c *gin.Context) {
// 	tokenString, err := c.Cookie("token")
// 	fmt.Println(err)
// 	if tokenString != "" && err == nil {
// 		if user, err := getUser(c); err != nil {
// 			model.UserToRedis(user)
// 		}
// 	}
// 	c.Next()
// }
