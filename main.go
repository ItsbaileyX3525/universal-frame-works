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

			var db *gorm.DB
			var DbErr error
			db, DbErr = connectDB(dbUsersName)
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

			fmt.Printf("User id: %s", userID)

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
				"localhost",
				false,
				true,
			)
			c.JSON(200, gin.H{
				"status": "Account Created!",
			})
		})

		api.POST("/login", func(c *gin.Context) {

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
