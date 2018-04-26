package test

import (
	"github.com/gin-gonic/gin"
)

// Hello 测试
func Hello(c *gin.Context) {
	// fmt.Println("\"message\":\"test\"")
	c.JSON(200, gin.H{"message": "test"})
}
