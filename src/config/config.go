package config

import (
	"fmt"
	"log"
	"os"

	"github.com/your-username/your-project/src/constants"
)

// Config holds the configuration for the application
type Config struct {
	ServerAddress string
	ConsumerHost  string
}

// LoadConfig loads configuration from environment variables or provides defaults
func LoadConfig() *Config {
	pricingServerHost := os.Getenv(constants.PricingServiceHost)
	fmt.Println(pricingServerHost)
	if pricingServerHost == "" {
		pricingServerHost = ":8000"
		log.Printf("SERVER_ADDRESS not set, using default '%s'", pricingServerHost)
	}
	consumerApiHost := os.Getenv(constants.ConsumerApi) 
	if consumerApiHost == "" {
		consumerApiHost = "https://consumer-api.development.dev.woltapi.com/home-assignment-api/" // fallback or blank
		log.Printf("CONSUMER_API not set, defaulting to %s", consumerApiHost)
	}

	return &Config{
		ServerAddress: pricingServerHost,
		ConsumerHost:  consumerApiHost,
	}
}
