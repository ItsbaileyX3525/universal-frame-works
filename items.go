package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handleGetItems(c *gin.Context) {
	var category string = c.Query("category")
	var page int
	var parseErr error

	page, parseErr = strconv.Atoi(c.Query("page"))
	if parseErr != nil {
		c.JSON(200, gin.H{"status": "error"})
		return
	}

	log.Print(category)
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

	type Item struct {
		ID        string  `json:"ID"`
		Name      string  `json:"Name"`
		ImagePath string  `json:"ImagePath"`
		AvgRating float32 `json:"AvgRating"`
	}

	var items []Item
	var err error

	err = db.Raw(
		"SELECT i.id, i.name, i.imagePath, COALESCE(AVG(r.rating), 0) as AvgRating FROM items i LEFT JOIN ratings r ON i.id = r.itemID WHERE i.category = ? GROUP BY i.id, i.name, i.imagePath, i.category",
		category,
	).Scan(&items).Error

	if err != nil {
		c.JSON(200, gin.H{"status": "db error", "error": err})
		return
	}

	if items == nil {
		items = []Item{}
	}

	c.JSON(200, gin.H{
		"status": "success",
		"items":  items,
	})
}

func handleGetAverageRating(c *gin.Context) {
	var uuid string = c.Query("uuid")

	var db *gorm.DB
	var DbErr error
	db, DbErr = connectDB()
	if DbErr != nil {
		c.JSON(200, gin.H{
			"status": "Something went wrong: " + DbErr.Error(),
		})
		return
	}

	var avgRating float32
	var row *sql.Row = db.Raw(
		"SELECT COALESCE(AVG(rating), 0) FROM ratings WHERE itemID = ?",
		uuid,
	).Row()

	if scanErr := row.Scan(&avgRating); scanErr != nil {
		c.JSON(200, gin.H{"status": "error"})
		return
	}

	c.JSON(200, gin.H{
		"status":    "success",
		"avgRating": avgRating,
	})
}
