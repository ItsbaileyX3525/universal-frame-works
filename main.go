package main

import (
	"os"
	"path/filepath"

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

	router.Static("/assets", "./assets")

	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "GET" {
			path := filepath.Join("./public", c.Request.URL.Path)

			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				c.File(path)
				return
			}

			if filepath.Ext(path) == "" {
				htmlPath := path + ".html"
				if _, err := os.Stat(htmlPath); err == nil {
					c.File(htmlPath)
					return
				}
			}

			if info, err := os.Stat(path); err == nil && info.IsDir() {
				indexPath := filepath.Join(path, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					c.File(indexPath)
					return
				}
			}

			c.File("./public/404.html")
		}
	})

	router.Run(":8080")
}
