package api

import (
	"github.com/gin-gonic/gin"
	"grpc-demo/grpc/client"
	"log"
)

func Command(c *gin.Context) {
	//host := c.PostForm("host")
	cmd := c.PostForm("cmd")
	hosts := []string{"127.0.0.1"}
	r, err := client.HandleExec(hosts, cmd)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"response": r,
	})
}