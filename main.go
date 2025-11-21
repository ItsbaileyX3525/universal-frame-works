package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbName string
var dbUser string
var dbPass string

// var dbUsersName string
var websiteURL string

func connectDB() (*gorm.DB, error) {
	var dsn string = fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", dbUser, dbPass, dbName)
	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Stack overflow
func generateRandomToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func createEndpoints(router *gin.Engine) {
	var api *gin.RouterGroup = router.Group("/api")
	{
		api.POST("/submitMessage", func(c *gin.Context) {
			//
		})

		api.POST("/submitRating", func(c *gin.Context) {
			var body struct {
				UUID   string `json:"uuid"`
				Rating int    `json:"rating"`
			}
			var err error = c.BindJSON(&body)
			if err != nil {
				c.JSON(200, gin.H{
					"error":  "Invalid POST",
					"errMsg": err,
				})
				return
			}
			var sessionID string
			sessionID, err = c.Cookie("session_id")
			if err != nil {
				c.JSON(200, gin.H{
					"status": "unauthorised",
				})
				return
			}

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB()
			if DbErr != nil {
				c.JSON(200, gin.H{
					"status": fmt.Sprintf("Something went wrong: %s", DbErr),
				})
				return
			}

			var userID string

			var row *sql.Row = db.Raw(
				"SELECT userID FROM sessions WHERE token = ?",
				sessionID,
			).Row()

			if scanErr := row.Scan(&userID); scanErr != nil {
				c.JSON(200, gin.H{"error": "Session not found"})
				return
			}

			var row2 *sql.Row = db.Raw(
				"SELECT ID FROM ratings WHERE userID = ? AND itemID = ?",
				userID,
				body.UUID,
			).Row()

			var existingID string

			var accountCheckErr error = row2.Scan(&existingID)
			if accountCheckErr == nil {
				c.JSON(200, gin.H{"error": "Rating already exists"})
				return
			}

			db.Exec(
				"INSERT INTO ratings (userID, itemID, rating) VALUES (?, ?, ?)",
				userID,
				body.UUID,
				body.Rating,
			)

			c.JSON(200, gin.H{
				"status":  "success",
				"message": "Rating submitted!",
			})
		})

		api.GET("/items", func(c *gin.Context) {
			var category string = c.Query("category")
			var page int
			var parseErr error

			page, parseErr = strconv.Atoi(c.Query("page"))
			if parseErr != nil {
				c.JSON(200, gin.H{"status": "Errorsss"})
				return
			}

			log.Print(category)
			log.Print(page)

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB()
			if DbErr != nil {
				c.JSON(200, gin.H{
					"status": fmt.Sprintf("Something went wrong: %s", DbErr),
				})
				return
			}

			type Item struct {
				ID   string
				Name string
			}

			var limit int = 8
			var offset int = (page - 1) * limit

			var items []Item
			var err error

			err = db.Raw(
				"SELECT id, name FROM items WHERE category = ? LIMIT ? OFFSET ?",
				category,
				limit,
				offset,
			).Scan(&items).Error

			if err != nil {
				c.JSON(200, gin.H{"status": "db error", "error": err})
				return
			}

			c.JSON(200, gin.H{
				"status": "success",
				"items":  items,
			})
		})

		api.POST("/signup", func(c *gin.Context) {
			var body struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			var err error = c.BindJSON(&body)
			if err != nil {
				c.JSON(200, gin.H{"error": "Invalid POST"})
				return
			}

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB()
			if DbErr != nil {
				c.JSON(200, gin.H{
					"status": fmt.Sprintf("Something went wrong: %s", DbErr),
				})
				return
			}

			var pass string = body.Password
			var user string = body.Username

			var existingID int

			var row *sql.Row = db.Raw(
				"SELECT id FROM users WHERE username = ? LIMIT 1",
				user,
			).Row()

			var accountCheckErr error = row.Scan(&existingID)
			if accountCheckErr == nil {
				c.JSON(200, gin.H{"error": "Username in use"})
				return
			}

			//Hash
			var bytes []byte
			var pwdErr error
			bytes, pwdErr = bcrypt.GenerateFromPassword([]byte(pass), 14)
			if pwdErr != nil {
				c.JSON(500, gin.H{"error": "Password hash failed"})
				return
			}
			//I think this is proper hashing? - Looks like it

			var sanitisedUser string = html.EscapeString(user)
			var encryptedPass string = string(bytes)

			var result *gorm.DB = db.Exec(
				"INSERT INTO users (username, password) VALUES (?, ?)",
				sanitisedUser,
				encryptedPass,
			)

			if result.Error != nil {
				c.JSON(500, gin.H{"error": "DB failed or something"})
				return
			}

			var userID string

			var row2 *sql.Row = db.Raw(
				"SELECT id FROM users WHERE username = ? LIMIT 1",
				user,
			).Row()

			var idRetrieve error = row2.Scan(&userID)
			if idRetrieve != nil {
				c.JSON(500, gin.H{"error": "Idek how this happened"})
				return
			}

			var sessionID string
			var sessionError error
			sessionID, sessionError = generateRandomToken()

			if sessionError != nil {
				c.JSON(500, gin.H{"error": "session error"})
				return
			}

			result = db.Exec(
				"INSERT INTO sessions (ID, userID, token, expiresAt) VALUES (?, ?, ?, ?)",
				sessionID,
				userID,
				sessionID,
				time.Now().Add(time.Hour*24*30),
			)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Failed to create session"})
				return
			}

			c.SetCookie(
				"session_id",
				sessionID,
				60*60*24*30,
				"/",
				websiteURL,
				false,
				true,
			)
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "Account Created!",
			})
		})

		api.POST("/login", func(c *gin.Context) {
			var body struct {
				Username        string `json:"username"`
				Password        string `json:"password"`
				ConfirmPassword string `json:"confirmPassword"`
			}

			var err error = c.BindJSON(&body)
			if err != nil {
				c.JSON(200, gin.H{"error": "Invalid POST"})
				return
			}

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB()
			if DbErr != nil {
				c.JSON(200, gin.H{
					"status": fmt.Sprintf("Something went wrong: %s", DbErr),
				})
				return
			}

			var username string = body.Username
			var password string = body.Password
			var confirmPassword string = body.ConfirmPassword

			if password != confirmPassword {
				c.JSON(200, gin.H{"error": "Password doesn't match"})
				return
			}

			var accountDetails struct {
				userID   int
				password string
			}

			var row *sql.Row = db.Raw(
				"SELECT id, password FROM users WHERE username = ? LIMIT 1",
				username,
			).Row()

			var accountCheckErr error = row.Scan(&accountDetails.userID, &accountDetails.password)
			if accountCheckErr != nil {
				c.JSON(200, gin.H{"error": "Account doesn't exist."})
				return
			}

			var pwdErr error = bcrypt.CompareHashAndPassword([]byte(accountDetails.password), []byte(password))
			if pwdErr != nil {
				c.JSON(200, gin.H{"error": "Invalid password"})
				return
			}

			var sessionResult *gorm.DB = db.Exec( //Ensures no stale tokens exist
				"DELETE FROM sessions WHERE userID = ?",
				accountDetails.userID,
			)

			if sessionResult.Error != nil {
				c.JSON(500, gin.H{"error": "Something went wrong deleting the session"})
				return
			}

			var sessionID string
			var sessionError error
			sessionID, sessionError = generateRandomToken()

			if sessionError != nil {
				c.JSON(500, gin.H{"error": "session error"})
				return
			}

			var result *gorm.DB = db.Exec(
				"INSERT INTO sessions (ID, userID, token, expiresAt) VALUES (?, ?, ?, ?)",
				sessionID,
				accountDetails.userID,
				sessionID,
				time.Now().Add(time.Hour*24*30),
			)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Failed to create session"})
				return
			}

			c.SetCookie(
				"session_id",
				sessionID,
				60*60*24*30,
				"/",
				websiteURL,
				false, //Change to true (I think it means https or something)
				true,
			)
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "Login successful!",
			})
		})

		api.POST("/requireLogin", func(c *gin.Context) {
			var sessionID string
			var err error
			sessionID, err = c.Cookie("session_id")
			log.Printf("Sesion: %s", sessionID)
			if err != nil {
				c.JSON(200, gin.H{
					"status": "unauthorised",
				})
				return
			}

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB()
			if DbErr != nil {
				c.JSON(500, gin.H{
					"status": fmt.Sprintf("Something went wrong: %s", DbErr),
				})
				return
			}

			row := db.Raw(
				"SELECT userID FROM sessions WHERE ID = ?",
				sessionID,
			).Row()
			var userID int
			if scanErr := row.Scan(&userID); scanErr != nil {
				c.JSON(200, gin.H{
					"status": "unauthorised",
				})
				return
			}

			c.JSON(200, gin.H{
				"status": "authenticated",
				"userID": userID,
			})
		})

		api.POST("/logout", func(c *gin.Context) {
			var sessionID string
			var err error
			sessionID, err = c.Cookie("session_id")
			if err != nil {
				c.AbortWithStatus(401)
				return
			}

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB()
			if DbErr != nil {
				c.JSON(200, gin.H{
					"status": fmt.Sprintf("Something went wrong: %s", DbErr),
				})
				return
			}

			c.SetCookie(
				"session_id",
				"",
				-1,
				"/",
				websiteURL,
				false,
				true,
			)

			db.Exec("DELETE FROM sessions WHERE ID=?", sessionID)
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "Logged out successfully.",
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
	websiteURL = os.Getenv("WEBSITE_NAME")

	var router *gin.Engine = gin.Default()
	//Remove
	//router.Use(func(c *gin.Context) {
	//c.Writer.Header().Set("Cache-Control", "no-cache")
	//})

	createEndpoints(router)
	renderHTML(router)

	router.Run(":8080")
}
