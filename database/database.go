package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// db is a package-level variable to hold the database connection
var db *gorm.DB

// Connect initializes the database connection with the provided parameters
func Connect(dbname, port, host, user, password, sslmode string) error {
	// Create the Data Source Name (DSN)
	var dsn string
	if password == "" {
		dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s",
			host, user, dbname, port, sslmode)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode)
	}

	var err error
	//db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable detailed logging
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

// GetDB returns the current database connection
func GetDB() *gorm.DB {
	return db
}
