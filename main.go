package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var dbName string
var dbUser string
var dbPass string
var websiteURL string
var cookieSecure bool

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
	var useTLS bool = true

	cookieSecureStr := os.Getenv("COOKIE_SECURE")
	cookieSecure, _ = strconv.ParseBool(cookieSecureStr)

	var router *gin.Engine = gin.Default()

	createEndpoints(router)
	renderHTML(router)

	_, fileTest := os.Stat("ssl/cert.pem")

	if errors.Is(fileTest, os.ErrNotExist) {
		useTLS = false
	}

	if useTLS {
		router.RunTLS(fmt.Sprintf(":%s", PORT), "ssl/cert.pem", "ssl/key.pem")
	} else {
		router.Run(fmt.Sprintf(":%s", PORT))
	}
}
