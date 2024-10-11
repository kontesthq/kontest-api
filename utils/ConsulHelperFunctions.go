package utils

import (
	"github.com/hashicorp/consul/api"
	"log"
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

func (c *ConsulService) Start(port int, serviceID string) {
	c.registerService(port, serviceID)
	go c.updateHealthCheck()
}

const (
	ttl     = time.Second * 10
	checkID = "checkAlive"
)

func (c *ConsulService) registerService(port int, serviceID string) {
	// Configure Consul client with the correct address and port
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:5150"

	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        checkID,
	}

	// Register a service with Connect sidecar proxy enabled
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "my-cluster",
		Tags:    []string{"kontest-api"},
		Address: "127.0.0.1",
		Port:    port,
		Check:   check,
		//Connect: &api.AgentServiceConnect{ // Enable Consul Connect (Service Mesh)
		//	SidecarService: &api.AgentServiceConnectProxyConfig{},
		//},
	}

	err := c.consulClient.Agent().ServiceRegister(serviceRegistration)
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

func (c *ConsulService) updateHealthCheck() {
	ticker := time.NewTicker(ttl / 2)

	for {
		err := c.consulClient.Agent().UpdateTTL(checkID, "Still alive", api.HealthPassing)

		if err != nil {
			log.Fatalf("Failed to update TTL: %v", err)
		}

		<-ticker.C
	}
}
