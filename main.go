package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbName string
var dbUser string
var dbPass string
var dbUsersName string

func connectDB(tblName string) (*gorm.DB, error) {
	var dsn string = fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", dbUser, dbPass, dbName)
	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createEndpoints(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/signup", func(c *gin.Context) {
			var body struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			var err error = c.BindJSON(&body)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid POST"})
				return
			}

			log.Printf("Username: %v", body.Username)
			log.Printf("Password: %v", body.Password)

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB(dbUsersName)
			if DbErr != nil {
				c.JSON(200, gin.H{
					"status": fmt.Sprintf("Something went wrong: %s", DbErr),
				})
				return
			}

			//Raw sql here

			c.JSON(200, gin.H{
				"status": "uhhh",
			})
		})

		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
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

	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbUsersName = os.Getenv("DB_USERS_NAME")

	log.Printf("Database username: %v", dbUser)
	log.Printf("Database password: %v", dbPass)

	var router *gin.Engine = gin.Default()

	createEndpoints(router)
	//connectDB()
	renderHTML(router)

	router.Run(":8080")
}
