package main

import "github.com/gin-gonic/gin"

func createEndpoints(router *gin.Engine) {
	var api *gin.RouterGroup = router.Group("/api")
	{
		api.GET("/comments", handleGetComments)
		api.POST("/submitMessage", handleSubmitMessage)

		api.POST("/submitRating", handleSubmitRating)
		api.GET("/averageRating", handleGetAverageRating)

		api.GET("/items", handleGetItems)

		api.POST("/signup", handleSignup)
		api.POST("/login", handleLogin)
		api.POST("/requireLogin", handleRequireLogin)
		api.POST("/logout", handleLogout)

		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
