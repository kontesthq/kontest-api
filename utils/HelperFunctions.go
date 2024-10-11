package utils

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"kontest-api/database"
	"log"
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

func RegisterService(port int, serviceID string) {
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

func DeregisterService(serviceID string) {
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
