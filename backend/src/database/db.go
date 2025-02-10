package database

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func init() {
	once.Do(func() {
		err := godotenv.Load("db.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		host := os.Getenv("DB_HOST")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		databaseName := os.Getenv("POSTGRES_DB")
		port := os.Getenv("DB_PORT")

		dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + databaseName + " port=" + port
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
	})
}
