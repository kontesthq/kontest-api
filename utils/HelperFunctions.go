package utils

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"kontest-api/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	if dbErr := database.Connect(dbname, dbPort, dbHost, user, password, sslmode); dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		return
	}

	// Create the uuid-ossp extension
	if err := database.GetDB().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create uuid-ossp extension: %v\n", err)
		return
	}

	createTables(database.GetDB())
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

func HandleShutdown(server *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting down server...")
		// Create a context with a timeout to allow for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %+v", err)
		}
		log.Println("Server exited properly")
	}()
}
