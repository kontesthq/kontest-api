package main

import (
	"fmt"
	"kontest-api/middleware"
	"kontest-api/routes"
	"kontest-api/utils"
	"log"
	"net/http"
	"os"
	"strconv"
)

var serviceName = "KONTEST-API"

func main() {
	port := os.Getenv("KONTEST_API_SERVER_PORT")
	if port == "" {
		port = "5151"
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Failed to convert port to integer: %v", err)
	}

	consulService := utils.NewConsulService()
	consulService.Start(portInt, serviceName)

	utils.InitalizeDatabase("kontest", "5432", "localhost", "postgres", "postgres", "disable")

	utils.InitializeDependencies()

	router := http.NewServeMux()

	routes.RegisterRoutes(router)

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":" + port,    // Use the field name Addr for the address
		Handler: stack(router), // Use the field name Handler for the router
	}

	fmt.Println("Server listening at port: " + port)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
