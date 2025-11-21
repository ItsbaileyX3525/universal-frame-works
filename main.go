package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var dbName string
var dbUser string
var dbPass string
var websiteURL string

func main() {
	var dotenvErr error = godotenv.Load()
	if dotenvErr != nil {
		log.Fatal(".env failed to load")
	}

	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	websiteURL = os.Getenv("WEBSITE_NAME")
	var PORT string = os.Getenv("PORT")

	var router *gin.Engine = gin.Default()

	createEndpoints(router)
	renderHTML(router)

	router.Run(fmt.Sprintf(":%s", PORT))
}
