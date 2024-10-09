package main

import (
	"fmt"
	"kontest-api/database"
	"kontest-api/middleware"
	"kontest-api/routes"
	"net/http"
	"os"
)

func main() {
	// Initialize the database connection parameters
	dbname := "kontest"
	dbPort := "5432"
	dbHost := "localhost"
	user := "ayushsinghal"
	password := ""
	sslmode := "disable"

	// Connect to the database
	if dbErr := database.Connect(dbname, dbPort, dbHost, user, password, sslmode); dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		return
	}

	router := http.NewServeMux()

	routes.RegisterRoutes(router)

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port,    // Use the field name Addr for the address
		Handler: stack(router), // Use the field name Handler for the router
	}

	fmt.Println("Server listening at port: " + port)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
