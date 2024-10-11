package utils

import (
	"fmt"
	"kontest-api/database"
	"os"
)

// InitalizeDatabase Initialize the database connection with default values
func InitalizeDatabase(
	dbname string,
	dbPort string,
	dbHost string,
	user string,
	password string,
	sslmode string,
) {

	if dbErr := database.Connect(dbname, dbPort, dbHost, user, password, sslmode); dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		return
	}
}
