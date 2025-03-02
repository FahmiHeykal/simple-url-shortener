package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "1234"
	DB_NAME     = "url_shortener"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

var DB *sql.DB

type URL struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func ConnectDB() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak bisa diakses:", err)
	}

	fmt.Println("Koneksi database sukses!")
}

func generateShortURL() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

func CreateShortURL(c *gin.Context) {
	var url URL

	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	url.ShortURL = generateShortURL()

	_, err := DB.Exec("INSERT INTO urls (original_url, short_url) VALUES ($1, $2)", url.OriginalURL, url.ShortURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original_url": url.OriginalURL,
		"short_url":    "http://localhost:8080/" + url.ShortURL,
	})
}

func RedirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")

	var originalURL string
	err := DB.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the URL Shortener API!"})
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/shorten", CreateShortURL)
	router.GET("/:shortURL", RedirectURL)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the URL Shortener API!",
		})
	})

	return router
}

func main() {
	ConnectDB()
	r := SetupRouter()

	fmt.Println("Server berjalan di http://localhost:8080")
	r.Run("127.0.0.1:8080")
}
