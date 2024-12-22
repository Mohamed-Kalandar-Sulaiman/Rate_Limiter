package utils

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// RateLimitConfig represents the rate limit configuration for a specific plan.
type RateLimitConfig struct {
	Name           string `yaml:"name"`
	Unit           string    `yaml:"unit"`
	RequestPerUnit int    `yaml:"request_per_unit"`
	Algorithm      string `yaml:"algorithm"`
}

// ActionConfig represents a specific action within a service.
type ActionConfig struct {
	Name      string           `yaml:"name"`
	ID        int              `yaml:"id"`
	RateLimit []RateLimitConfig `yaml:"rate_limit"`
}

// ServiceActionConfig holds the configuration data for a service with actions.
type ServiceActionConfig struct {
	ServiceName string         `yaml:"name"`
	ServiceID   int            `yaml:"id"`
	Actions     []ActionConfig `yaml:"actions"`
}

// ServicesConfig wraps the list of service configurations from the YAML file.
type ServicesConfig struct {
	Services []ServiceActionConfig `yaml:"services"`
}

// ConfigData represents the flattened configuration details.
type ConfigData struct {
	ServiceID       int             `yaml:"service_id"`
	ActionID        int             `yaml:"action_id"`
	ConfigID        string          `yaml:"config_id"`
	RateLimitConfig RateLimitConfig `yaml:"rate_limit"`
}

// ConfigMap holds the flattened in-memory configuration.
type ConfigMap struct {
	Config       map[string]ConfigData
	LastRefreshed time.Time
}

// New creates and returns a new instance of ConfigMap with an empty configuration.
func NewConfigMap() *ConfigMap {
	return &ConfigMap{
		Config: make(map[string]ConfigData),
	}
}

// LoadConfig loads the configuration from the YAML file and flattens it.
func (c *ConfigMap) LoadConfig(filename string) error {
	// Read the YAML file using os.ReadFile
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into ServicesConfig
	var servicesConfig ServicesConfig
	err = yaml.Unmarshal(data, &servicesConfig)
	if err != nil {
		return err
	}

	// Flatten the structure into ConfigMap
	c.Config = make(map[string]ConfigData)
	for _, service := range servicesConfig.Services {
		for _, action := range service.Actions {
			for _, rateLimit := range action.RateLimit {
				key := fmt.Sprintf("%s:%s:%s", service.ServiceName, action.Name, rateLimit.Name)
				c.Config[key] = ConfigData{
					ServiceID:       service.ServiceID,
					ActionID:        action.ID,
					ConfigID:        rateLimit.Name,
					RateLimitConfig: rateLimit,
				}
			}
		}
	}

	// Set last refreshed time
	c.LastRefreshed = time.Now()

	return nil
}

// GetConfig fetches the configuration based on the key. If more than 1 hour has passed since the last refresh, it reloads the config.
func (c *ConfigMap) GetConfig(key string) (*ConfigData, error) {
	// If more than 1 hour has passed since the last refresh, reload the config
	if time.Since(c.LastRefreshed) > time.Minute {
		fmt.Println("Refreshing config...")
		// Replace this with the actual file path you want to load from
		err := c.LoadConfig("config.yaml")
		if err != nil {
			return nil, err
		}
	}

	// Return the config for the given key
	configData, exists := c.Config[key]
	if !exists {
		return nil, fmt.Errorf("configuration not found for key: %s", key)
	}

	return &configData, nil
}
