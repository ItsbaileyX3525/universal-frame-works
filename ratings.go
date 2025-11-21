package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleSubmitRating(c *gin.Context) {
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

	var row2 *sql.Row = db.Raw(
		"SELECT ID FROM ratings WHERE userID = ? AND itemID = ?",
		userID,
		body.UUID,
	).Row()

	var existingID string
	var ratingExists bool = false

	var accountCheckErr error = row2.Scan(&existingID)
	if accountCheckErr == nil {
		ratingExists = true
	}

	if ratingExists {
		db.Exec(
			"UPDATE ratings SET rating = ? WHERE itemID = ? AND userID = ?",
			body.Rating,
			body.UUID,
			userID,
		)
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Rating updated!",
		})
	} else {
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
	}
}
