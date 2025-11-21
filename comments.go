package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleGetComments(c *gin.Context) {
	var uuid string = c.Query("uuid")
	var page int
	var parseErr error

	page, parseErr = strconv.Atoi(c.Query("page"))
	if parseErr != nil {
		c.JSON(200, gin.H{"status": "error"})
		return
	}

	log.Print(uuid)
	log.Print(page)

	var db *gorm.DB
	var DbErr error
	db, DbErr = connectDB()
	if DbErr != nil {
		c.JSON(200, gin.H{
			"status": "Something went wrong: " + DbErr.Error(),
		})
		return
	}

	type Comment struct {
		Username string `json:"Username"`
		Content  string `json:"Comment"`
		Creation string `json:"creation"`
	}

	var comments []Comment
	var err error

	err = db.Raw(
		"SELECT u.username as Username, c.content as Content, c.creation FROM comments c JOIN users u ON c.userID = u.ID WHERE c.itemID = ? LIMIT 5",
		uuid,
	).Scan(&comments).Error

	if err != nil {
		c.JSON(200, gin.H{"status": "db error", "error": err})
		return
	}

	if comments == nil {
		comments = []Comment{}
	}

	c.JSON(200, gin.H{
		"status":   "success",
		"comments": comments,
	})
}

func handleSubmitMessage(c *gin.Context) {
	var body struct {
		UUID    string `json:"uuid"`
		Message string `json:"message"`
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
			"status": "Something went wrong: " + DbErr.Error(),
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

	var result *gorm.DB = db.Exec(
		"INSERT INTO comments (userID, itemID, content) VALUES (?, ?, ?)",
		userID,
		body.UUID,
		body.Message,
	)

	if result.Error != nil {
		c.JSON(200, gin.H{"status": "error", "message": "Failed to submit comment"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Comment submitted!",
	})
}
