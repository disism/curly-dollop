package main

import (
	"github.com/gin-gonic/gin"
	"grpc-demo/api"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/command", api.Command)
	r.Run(":8082")
}