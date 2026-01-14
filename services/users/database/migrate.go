package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/tagaertner/e-commerce-graphql/services/users/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	// Try DATABASE_URL first for Render/production
	dbURL := os.Getenv("DATABASE_URL")

	// Fallback to individual vars for local deb
	if dbURL == ""{
		dbURL = fmt.Sprintf(
	
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("DB_PORT"),
	)

		log.Println("🔧 Using individual DB environment variables")
	} else {
		log.Println("🔗 Using DATABASE_URL connection string")
	}

	maxRetries := 20
	retryDelay := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			log.Println("✅ Connected to PostgreSQL successfully")
			return db
		}

		log.Printf("❌ Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			log.Printf("⏳ Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	log.Fatalf("❌ Could not connect to database after %d retries", maxRetries)
	return nil
}

func RunMigrations(db *gorm.DB) {
    log.Println("📦 Running AutoMigrate...")
    if err := db.AutoMigrate(&models.User{}); err != nil {
        log.Fatalf("❌ Migration failed: %v", err)
    }
    log.Println("✅ Migrations complete")
}