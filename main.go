package main

import (
	"fmt"
	consulServiceManager "github.com/ayushs-2k4/go-consul-service-manager"
	loadbalancer "github.com/ayushs-2k4/go-load-balancer"
	"io"
	"kontest-api/database"
	"kontest-api/middleware"
	"kontest-api/routes"
	"kontest-api/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	applicationHost = "localhost"   // Default value for local development
	applicationPort = "5151"        // Default value for local development
	serviceName     = "KONTEST-API" // Service name for Service Registry
	consulHost      = "localhost"   // Default value for local development
	consulPort      = 5150          // Port as a constant (can be constant if it won't change)

	dbHost           = "localhost"
	dbPort           = "5432"
	dbName           = "kontest"
	dbUser           = "ayushsinghal"
	dbPassword       = ""
	isSSLModeEnabled = false
)

func initializeVariables() {
	// Attempt to read the KONTEST_API_SERVER_HOST environment variable
	if host := os.Getenv("KONTEST_API_SERVER_HOST"); host != "" {
		applicationHost = host // Override with the environment variable if set
	}

	// Attempt to read the KONTEST_API_SERVER_PORT environment variable
	if port := os.Getenv("KONTEST_API_SERVER_PORT"); port != "" {
		applicationPort = port // Override with the environment variable if set
	}

	// Attempt to read the CONSUL_ADDRESS environment variable
	if host := os.Getenv("CONSUL_HOST"); host != "" {
		consulHost = host // Override with the environment variable if set
	}

	// Attempt to read the CONSUL_PORT environment variable
	if port := os.Getenv("CONSUL_PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil {
			consulPort = portInt // Override with the environment variable if set and valid
		}
	}

	// Attempt to read the DB_HOST environment variable
	if host := os.Getenv("DB_HOST"); host != "" {
		dbHost = host // Override with the environment variable if set
	}

	// Attempt to read the DB_PORT environment variable
	if port := os.Getenv("DB_PORT"); port != "" {
		dbPort = port // Override with the environment variable if set
	}

	// Attempt to read the DB_NAME environment variable
	if name := os.Getenv("DB_NAME"); name != "" {
		dbName = name // Override with the environment variable if set
	}

	// Attempt to read the DB_USER environment variable
	if user := os.Getenv("DB_USER"); user != "" {
		dbUser = user // Override with the environment variable if set
	}

	// Attempt to read the DB_PASSWORD environment variable
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		dbPassword = password // Override with the environment variable if set
	}

	// Attempt to read the DB_SSL_MODE environment variable
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		isSSLModeEnabled = sslMode == "enable"
	}
}

func main() {
	initializeVariables()

	portInt, err := strconv.Atoi(applicationPort)
	if err != nil {
		log.Fatalf("Failed to convert applicationPort to integer: %v", err)
	}

	consulService := consulServiceManager.NewConsulService(consulHost, consulPort)
	consulService.Start(applicationHost, portInt, serviceName)

	//checkingLoadBalancer()
	//checkLoadBalancerUserStatsService()

	database.InitalizeDatabase(dbPort, dbHost, dbUser, dbPassword, dbName, map[bool]string{true: "enable", false: "disable"}[isSSLModeEnabled])

	utils.InitializeDependencies()

	router := http.NewServeMux()

	routes.RegisterRoutes(router)

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":" + applicationPort, // Use the field name Addr for the address
		Handler: stack(router),         // Use the field name Handler for the router
	}

	fmt.Println("Server listening at applicationPort: " + applicationPort)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func checkingLoadBalancer() {
	lb, err := loadbalancer.GetLoadBalancer(serviceName, consulHost, consulPort)

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

func checkLoadBalancerUserStatsService() {
	lb, err := loadbalancer.GetLoadBalancer("KONTEST-USER-STATS-SERVICE", consulHost, consulPort)

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
	url := fmt.Sprintf("http://%s:%d/codechef_user?username=ayushs_2k4", instance.Address, instance.Port)
	fmt.Printf("Calling URL: %s\n", url)

	// Make the HTTP GET request to the service
	client := &http.Client{
		Timeout: 5 * time.Second, // Set a timeout to avoid hanging requests
	}

	secondsToWait := 2

	// Wait for some time without blocking the main thread
	go func() {
		fmt.Println("Waiting for" + strconv.Itoa(secondsToWait) + " seconds...")
		time.Sleep(time.Duration(secondsToWait) * time.Second)
		fmt.Println(strconv.Itoa(secondsToWait) + " seconds passed, making HTTP request.")

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
