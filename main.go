package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func generateAlphanumericID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type Url struct {
	gorm.Model
	UriID string
	Uri   string
}

type UrlShortenRequest struct {
	Url string
}

var DB *gorm.DB

func pingEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func shorten(c *gin.Context) {
	var requestBody UrlShortenRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusInternalServerError, "")
	}

	insertPayload := &Url{UriID: generateAlphanumericID(8), Uri: requestBody.Url}

	DB.Model(&Url{}).Create(insertPayload)

	c.JSON(http.StatusOK, gin.H{
		"url": os.Getenv("BASE_URL") + insertPayload.UriID,
	})
}

func matchUrl(c *gin.Context) {
	UriID := c.Param("UriID")
	var result Url
	DB.Model(&Url{UriID: UriID}).First(&result)

	http.Redirect(c.Writer, c.Request, result.Uri, 303)
}

func main() {
	godotenv.Load()
	r := gin.Default()
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_PATH")), &gorm.Config{})

	db.AutoMigrate(&Url{})

	if err != nil {
		log.Fatal(err)
	}

	DB = db
	r.GET("/ping", pingEndpoint)
	r.POST("/shorten", shorten)
	r.GET("/:UriID", matchUrl)

	r.Run()
}
