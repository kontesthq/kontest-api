package database

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
)

// InitalizeDatabase Initialize the database connection with default values
func InitalizeDatabase(
	dbPort string,
	dbHost string,
	user string,
	password string,
	dbname string,
	sslmode string,
) {

	if dbErr := Connect(dbname, dbPort, dbHost, user, password, sslmode); dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		return
	}

	// Create the uuid-ossp extension
	if err := GetDB().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create uuid-ossp extension: %v\n", err)
		return
	}

	createTables(GetDB())
}

// createTables creates the necessary tables in the database
func createTables(db *gorm.DB) {
	// SQL query to create the 'kontests' table
	createKontestsTable := `
    CREATE TABLE IF NOT EXISTS kontests (
        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
        name TEXT NOT NULL,
        url TEXT,
        start_time TEXT,
        end_time TEXT,
        location TEXT,
        status TEXT,
        site_abbreviation TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	// Execute the SQL query for the 'kontests' table
	if err := db.Exec(createKontestsTable).Error; err != nil {
		log.Fatalf("Error creating kontests table: %v", err)
	}

	// SQL query to create the 'kontests_metadata' table
	createMetadataTable := `
    CREATE TABLE IF NOT EXISTS kontests_metadata (
        id TEXT PRIMARY KEY,
        last_updated_at TIMESTAMP
    );`

	// Execute the SQL query for the 'kontests_metadata' table
	if err := db.Exec(createMetadataTable).Error; err != nil {
		log.Fatalf("Error creating kontests_metadata table: %v", err)
	}
}
