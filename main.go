package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbUser string
var dbPass string

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

func connectDB(usr string, pwd string, dbName string) (*gorm.DB, error) {
	var dsn string = usr + ":" + pwd + "@tcp(localhost:3306)/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func renderHTML(router *gin.Engine) {
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
}

func main() {
	var dotenvErr error = godotenv.Load()
	if dotenvErr != nil {
		log.Fatal(".env failed to load")
	}

	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")

	log.Printf("Database username: %v", dbUser)
	log.Printf("Database password: %v", dbPass)

	var router *gin.Engine = gin.Default()

	createEndpoints(router)
	//connectDB()
	renderHTML(router)

	router.Run(":8080")
}
