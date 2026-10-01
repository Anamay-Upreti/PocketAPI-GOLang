package database

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"pocketapi/internal/config"
	"pocketapi/internal/models"
)

var DB *gorm.DB

// ConnectDatabase initializes PostgreSQL connection.
func ConnectDatabase() {

	var err error

	DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	sqlDB, err := DB.DB()

	if err != nil {
		log.Fatal(err)
	}

	configureConnectionPool(sqlDB)

	autoMigrate()

	log.Println("PostgreSQL Connected Successfully")
}

func autoMigrate() {

	err := DB.AutoMigrate(
	&models.User{},
	&models.Note{},
	&models.RefreshToken{},
)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database Migration Completed")
}

func configureConnectionPool(sqlDB *sql.DB) {

	sqlDB.SetMaxIdleConns(5)

	sqlDB.SetMaxOpenConns(20)

	sqlDB.SetConnMaxLifetime(time.Hour)

	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
}
