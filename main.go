package main

import (
	"fmt"
	loadbalancer "github.com/ayushs-2k4/go-load-balancer"
	"io"
	"kontest-api/middleware"
	"kontest-api/routes"
	"kontest-api/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

	consulService := utils.NewConsulService("localhost", 5150)
	consulService.Start(portInt, serviceName)

	//checkingLoadBalancer()

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

func checkingLoadBalancer() {
	lb, err := loadbalancer.GetLoadBalancer(serviceName, "localhost", 5150)

	if err != nil {
		log.Fatalf("Failed to create load balancer: %v", err)
	}

	// Choose a random instance
	instance, err := lb.ChooseInstance()
	if err != nil || instance == nil {
		log.Fatalf("Failed to choose instance: %v", err)
	}

	// Print the instance's address
	fmt.Printf("Instance address: %s:%d\n", instance.Address, instance.Port)

	// Construct the URL for the chosen instance
	url := fmt.Sprintf("http://%s:%d/kontests?page=1&per_page=10", instance.Address, instance.Port)
	fmt.Printf("Calling URL: %s\n", url)

	// Make the HTTP GET request to the service
	client := &http.Client{
		Timeout: 5 * time.Second, // Set a timeout to avoid hanging requests
	}

	// Wait for 5 seconds without blocking the main thread
	go func() {
		fmt.Println("Waiting for 5 seconds...")
		time.Sleep(5 * time.Second)
		fmt.Println("5 seconds passed, making HTTP request.")

		// Perform the HTTP request in the Goroutine
		resp, err := client.Get(url)
		if err != nil {
			log.Fatalf("Failed to make HTTP request: %v", err)
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response: %v", err)
		}

		// Print the response
		fmt.Printf("Response from instance: %s\n", body)
	}()

	// Continue doing other work in the main thread if needed
	fmt.Println("Main thread is not blocked, continuing execution...")
}
