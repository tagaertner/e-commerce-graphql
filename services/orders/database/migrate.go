package database

import (
	"fmt"
	"log"
	"os"
	"github.com/tagaertner/e-commerce-graphql/services/orders/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)
type Order struct {
	ID	string	`gorm:"primaryKey"`
	UserID string
	Quantity int
	TotalPrice float64
	Status string
	CreatAt time.Time
	Products []Product `gorm:"many2many:order_products;"`
}

type Product struct {
	ID string `gorm:"primaryKey"`
}
func Connect() *gorm.DB {

		dbURL := os.Getenv("DATABASE_URL")

	// Fallback to individual vars for local dev
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
	db.AutoMigrate(&models.Order{}, &models.Product{})
}