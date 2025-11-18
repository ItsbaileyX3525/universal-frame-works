package main

import (
	"github.com/gin-gonic/gin"
)

func createEndpoints(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "endpoint registered.",
			})
		})

		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}

func main() {
	router := gin.Default()
	createEndpoints(router)

	//router.Static("/", "./public") - Fix

	router.Run(":8080")
}
