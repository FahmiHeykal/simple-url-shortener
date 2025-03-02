package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres" // Ganti dengan username PostgreSQL kamu
	DB_PASSWORD = "1234"     // Ganti dengan password PostgreSQL kamu
	DB_NAME     = "url_shortener"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

var DB *sql.DB

func ConnectDB() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak dapat diakses:", err)
	}

	fmt.Println("✅ Koneksi database sukses!")
}
