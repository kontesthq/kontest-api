package utils

import (
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ConsulService struct {
	consulClient *api.Client
}

func NewConsulService() *ConsulService {
	// Configure Consul client with the correct address and port
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:5150"

	client, err := api.NewClient(config)

	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
		return nil
	}

	return &ConsulService{
		consulClient: client,
	}
}

func (c *ConsulService) Start(port int, serviceName string) {
	serviceID := serviceName + "-" + "1"

	// Ensure the service is deregistered when the application shuts down
	defer deregisterService(serviceID) // Will run when Start() exits

	// Register service and start health check
	c.registerService(port, serviceName, serviceID)
	go c.updateHealthCheck(serviceID)

	// Set up signal handling to wait for termination signals
	c.handleShutdown(serviceID)
}

const (
	ttl     = time.Second * 10
	checkID = "checkAlive"
)

func (c *ConsulService) registerService(port int, serviceName string, serviceID string) {
	// Configure Consul client with the correct address and port
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:5150"

	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        serviceID + "-" + checkID,
	}

	// Register a service with Connect sidecar proxy enabled
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Tags:    []string{"kontest-api"},
		Address: "127.0.0.1",
		Port:    port,
		Check:   check,
		Connect: &api.AgentServiceConnect{ // Enable Consul Connect (Service Mesh)
			Native: false,
			SidecarService: &api.AgentServiceRegistration{
				Name:    serviceName + "-sidecar",
				Tags:    []string{"sidecar"},
				Address: "127.0.0.1",
				Port:    20000,
				Check:   check,
			},
		},
	}

	err := c.consulClient.Agent().ServiceRegister(serviceRegistration)
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	log.Println("Service registered with Consul on port 5150")
}

func deregisterService(serviceID string) {
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

func (c *ConsulService) updateHealthCheck(serviceID string) {
	ticker := time.NewTicker(ttl / 2)
	finalID := serviceID + "-" + checkID

	for {
		err := c.consulClient.Agent().UpdateTTL(finalID, "Still alive", api.HealthPassing)

		if err != nil {
			log.Fatalf("Failed to update TTL: %v", err)
		}

		<-ticker.C
	}
}

// handleShutdown waits for termination signals and calls deregisterService
func (c *ConsulService) handleShutdown(serviceID string) {
	// Create a channel to listen for OS signals
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigChannel

	// Log and deregister the service
	log.Printf("Received signal: %s, deregistering service...", sig)
	deregisterService(serviceID)

	// Exit the application gracefully
	log.Println("Exiting application")
	os.Exit(0)
}
