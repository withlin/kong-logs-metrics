package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Aaaa(c *gin.Context) {
	fmt.Println("来啦")
	fmt.Println(c.Request.Header)
	c.Next()
}
