package main

import (
	"fmt"
	"github.com/ayushs-2k4/go-consul-service-manager/consulservicemanager"
	"github.com/ayushs-2k4/go-load-balancer/loadbalancer"
	"io"
	"kontest-api/database"
	"kontest-api/middleware"
	"kontest-api/routes"
	"kontest-api/utils"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	applicationHost = "localhost"   // Default value for local development
	applicationPort = 5151          // Default value for local development
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
	// Get the hostname of the machine
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Error fetching hostname", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Attempt to read the KONTEST_API_SERVER_HOST environment variable
	if host := os.Getenv("KONTEST_API_SERVER_HOST"); host != "" {
		applicationHost = host // Override with the environment variable if set
	} else {
		applicationHost = hostname // Use the machine's hostname if the env var is not set
	}

	// Attempt to read the KONTEST_API_SERVER_PORT environment variable
	if port := os.Getenv("KONTEST_API_SERVER_PORT"); port != "" {
		parsedPort, err := strconv.Atoi(port)
		if err != nil {
			slog.Error("Invalid port value", slog.String("error", err.Error()), slog.String("port", port))
			os.Exit(1) // Exit the program with a non-zero status code
		}
		applicationPort = parsedPort // Override with the environment variable if set
		slog.Info("Application port set from environment variable", slog.Int("applicationPort", applicationPort))
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

func setupLogging() *os.File {
	LOG_FILE := os.Getenv("LOG_FILE")

	if LOG_FILE == "" {
		LOG_FILE = "tmp/logs/logs.log"
	}

	// Get the directory from the log file path
	logDir := filepath.Dir(LOG_FILE)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		slog.Error("Failed to create log directory", slog.String("error", err.Error()))
		os.Exit(1)
	}

	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
	}
	// Open or create a log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// Handle error if the log file cannot be opened or created
		slog.Error("Failed to open log file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	w := io.MultiWriter(os.Stdout, logFile)

	// Configure slog to output JSON
	slog.SetDefault(slog.New(slog.NewJSONHandler(w, handlerOptions)))

	// Return the log file to close it in the main function
	return logFile
}

func main() {
	// Set up logging and capture the log file
	logFile := setupLogging()

	// Ensure the log file is closed when the program exits
	defer logFile.Close()

	// Log server restart with a timestamp
	slog.Info("Server restarted", slog.Time("time", time.Now()))

	initializeVariables()

	consulService := consulservicemanager.NewConsulService(consulHost, consulPort)
	consulService.Start(applicationHost, applicationPort, serviceName, []string{})

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
		Addr:    ":" + strconv.Itoa(applicationPort), // Use the field name Addr for the address
		Handler: stack(router),                       // Use the field name Handler for the router
	}

	slog.Info("Server listening", slog.Int("port", applicationPort))

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func checkingLoadBalancer() {
	lb, err := loadbalancer.GetLoadBalancer(serviceName, consulHost, consulPort)

	if err != nil {
		slog.Error("Failed to create load balancer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	instance, err := lb.ChooseInstance()
	if err != nil || instance == nil {
		slog.Error("Failed to choose instance", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Chosen instance", slog.String("address", instance.Address), slog.Int("port", instance.Port))

	url := fmt.Sprintf("http://%s:%d/kontests?page=1&per_page=10", instance.Address, instance.Port)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	go func() {
		time.Sleep(5 * time.Second)
		resp, err := client.Get(url)
		if err != nil {
			slog.Error("Failed to make HTTP request", slog.String("error", err.Error()))
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Failed to read response", slog.String("error", err.Error()))
			os.Exit(1)
		}

		slog.Info("Response from instance", slog.String("response", string(body)))
	}()

	slog.Info("Main thread is continuing execution...")
}

func checkLoadBalancerUserStatsService() {
	lb, err := loadbalancer.GetLoadBalancer("KONTEST-USER-STATS-SERVICE", consulHost, consulPort)

	if err != nil {
		slog.Error("Failed to create load balancer", slog.String("error", err.Error()))
		os.Exit(1) // Exit with an error
	}

	// Choose a random instance
	instance, err := lb.ChooseInstance()
	if err != nil || instance == nil {
		slog.Error("Failed to choose instance", slog.String("error", err.Error()))
		os.Exit(1) // Exit with an error
	}

	// Log the instance's address
	slog.Info("Instance chosen", slog.String("address", instance.Address), slog.Int("port", instance.Port))

	// Construct the URL for the chosen instance
	url := fmt.Sprintf("http://%s:%d/codechef_user?username=ayushs_2k4", instance.Address, instance.Port)
	slog.Info("Calling URL", slog.String("url", url))

	// Make the HTTP GET request to the service
	client := &http.Client{
		Timeout: 5 * time.Second, // Set a timeout to avoid hanging requests
	}

	secondsToWait := 2

	// Wait for some time without blocking the main thread
	go func() {
		slog.Info("Waiting before making the request", slog.Int("seconds", secondsToWait))
		time.Sleep(time.Duration(secondsToWait) * time.Second)
		slog.Info("Making HTTP request", slog.Int("seconds waited", secondsToWait))

		// Perform the HTTP request in the Goroutine
		resp, err := client.Get(url)
		if err != nil {
			slog.Error("Failed to make HTTP request", slog.String("error", err.Error()))
			os.Exit(1)
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Failed to read response", slog.String("error", err.Error()))
			os.Exit(1)
		}

		// Log the response
		slog.Info("Response from instance", slog.String("response", string(body)))
	}()

	// Log that the main thread is continuing execution
	slog.Info("Main thread is not blocked, continuing execution...")
}
