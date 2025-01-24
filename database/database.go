package database

import (
	"donate-bot/config"
	"donate-bot/models"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s",
		config.Config("POSTGRES_HOST"),
		config.Config("POSTGRES_USER"),
		config.Config("POSTGRES_PASSWORD"),
		config.Config("POSTGRES_PORT"),
		"postgres",
	)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	checkDBExists := `
		SELECT 1 FROM pg_database WHERE datname = $1;
	`
	var exists int
	err = DB.Raw(checkDBExists, config.Config("POSTGRES_DATABASE_NAME")).Scan(&exists).Error
	if err != nil {
		log.Fatal("Failed to check database existence:", err)
	}

	if exists == 0 {
		createDB := fmt.Sprintf("CREATE DATABASE %s;", config.Config("POSTGRES_DATABASE_NAME"))
		if err := DB.Exec(createDB).Error; err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Println("Database created successfully.")
	}

	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.Config("POSTGRES_HOST"),
		config.Config("POSTGRES_USER"),
		config.Config("POSTGRES_PASSWORD"),
		config.Config("POSTGRES_DATABASE_NAME"),
		config.Config("POSTGRES_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the target database:", err)
	}

	if err := DB.AutoMigrate(
		&models.DonationHistory{},
		&models.Referral{},
		&models.User{},
	); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	return DB
}
