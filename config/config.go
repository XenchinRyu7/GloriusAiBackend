package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"gloriusaiapi/models"

	"github.com/joho/godotenv"
)

var DB *gorm.DB

// InitializeDB initializes the database connection
func InitializeDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Debug: Print connection details (optional)
	fmt.Printf("Connecting to database: %s@tcp(%s:%s)/%s\n", dbUser, dbHost, dbPort, dbName)

	// Build connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database
	DB, err = gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Automigrate models
	err = DB.AutoMigrate(&models.User{}, &models.Message{}, &models.Context{}).Error
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connection established successfully!")
}
