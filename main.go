package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"kontest-api/database"
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

	registerService(portInt)
	defer deregisterService() // Ensure the service is deregistered when exiting

	initalizeDatabase("kontest", "5432", "localhost", "postgres", "postgres", "disable")

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

func registerService(port int) {
	// Configure Consul client with the correct address and port
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:5150" // Change to the correct Consul agent address and port

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
		return
	}

	// Register a service
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceID,
		Address: "127.0.0.1",
		Port:    port,
		Tags:    []string{"primary"},
	}

	err = client.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	log.Println("Service registered with Consul on port 5150")
}

func deregisterService() {
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:5150" // Ensure the correct address and port are used

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
		return
	}

	// Deregister the service
	err = client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		log.Fatalf("Failed to deregister service: %v", err)
	}

	log.Println(serviceID + " Service deregistered from Consul")
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
