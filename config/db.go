package config

import (
	"fmt"
	"log"
	"os"

	"github.com/herdiansc/go-cms/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dbOpen(port string) (*gorm.DB, error) {
	if port == "" {
		port = "5432"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), port)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func SetupDB(port string) *gorm.DB {
	DB, err := dbOpen(port)
	if err != nil {
		log.Fatalf("Error connecting to the database: %+v\n", err)
	}
	log.Printf("DB connection established: %+v\n", DB)
	DB.AutoMigrate(&models.Auth{})
	DB.AutoMigrate(&models.Article{})
	DB.AutoMigrate(&models.ArticleTag{})
	DB.AutoMigrate(&models.Tag{})
	DB.AutoMigrate(&models.TagTrendingScore{})
	DB.AutoMigrate(&models.ArticleHistory{})
	return DB
}
