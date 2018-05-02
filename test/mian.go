package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type CreateParams struct {
	Username     string   `json:"username"`
	Guests       []Person `json:"guests"`
	RoomType     string   `json:"roomType"`
	CheckinDate  string   `json:"checkinDate"`
	CheckoutDate string   `json:"checkoutDate"`
}

type Person struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		var createParams CreateParams
		err := c.BindJSON(&createParams)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(createParams)
	})

	r.Run(":8000")

	// request
	// curl http://localhost:8000 -d @request.json
}
