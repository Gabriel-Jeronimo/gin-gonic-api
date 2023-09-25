package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model
	id  string
	url string
}

var DB *gorm.DB

func pingEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func shorten(c *gin.Context) {
	
	if 
	c.JSON(http.StatusOK, gin.H{
		"": "",
	})
}

func main() {
	r := gin.Default()
	api := r.Group("/api")
	db, err := gorm.Open(sqlite.Open("url.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Url{})

	DB = db
	api.GET("/ping", pingEndpoint)
	api.POST("/shorten", shorten)

	r.Run()
}
