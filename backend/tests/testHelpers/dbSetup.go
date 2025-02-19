package testHelpers

import (
	"RemitlyTask/src/models"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	err := godotenv.Load("../../db.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DB_PORT")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + databaseName + " port=" + port
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.SwiftCode{})
	return db
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
	err := db.Migrator().DropTable(&models.SwiftCode{})
	if err != nil {
		t.Fatalf("Failed to clean up test database: %v", err)
	}
}
