package main

import (
	"fmt"
	"kontest-api/database"
	"kontest-api/middleware"
	"kontest-api/routes"
	"kontest-api/utils"
	"net/http"
	"os"
)

func main() {
	initalizeDatabase("kontest", "5432", "localhost", "postgres", "postgres", "disable")

	utils.InitializeDependencies()

	router := http.NewServeMux()

	routes.RegisterRoutes(router)

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	port := os.Getenv("KONTEST_API_SERVER_PORT")
	if port == "" {
		port = "5151"
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

// Initialize the database connection with default values
func initalizeDatabase(
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
