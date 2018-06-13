package middleware

import (
	"github.com/gin-gonic/gin"
)

//KeyAuth key验证
func KeyAuth(c *gin.Context) {
	if c.GetHeader("Autorizacion") != "a" {

		c.JSON(401, gin.H{"message": "没有权限访问"})
		return
	} else {
		c.Next()
	}

}
