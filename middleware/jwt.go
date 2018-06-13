package middleware

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

//JWT
func JWT(c *gin.Context) {
	result, err := ioutil.ReadAll(c.Request.Body)
	c.SetCookie("token", "aaaa", 60, "/", "", true, true)
	if err != nil {
		fmt.Println("=========================")
		fmt.Println(err.Error())
		fmt.Println("=========================")
	}
	fmt.Println("=========================")
	fmt.Println(string(result))
	fmt.Println("=========================")
	fmt.Println(c.Request.Cookie("aaaa"))
	c.Next()

}
