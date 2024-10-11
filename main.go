package main

import (
	"context"
	"fmt"
	"kontest-api/middleware"
	"kontest-api/routes"
	"kontest-api/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var serviceID = "KONTEST-API"

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
	consulService.Start(portInt, serviceID)

	defer utils.DeregisterService(serviceID) // Ensure the service is deregistered when exiting

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

	// Handle termination signals
	handleShutdown(&server)

	fmt.Println("Server listening at port: " + port)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleShutdown(server *http.Server) {
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
