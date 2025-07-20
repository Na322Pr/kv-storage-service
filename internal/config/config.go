package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Node `yaml:"node" env-required:"true"`
	GRPC `yaml:"grpc" env-required:"true"`
}

type Node struct {
	ID        int      `yaml:"id" env:"NODE_ID" env-required:"true"`
	SeedNodes []string `yaml:"seed_nodes" env:"SEED_NODES" env-separator:","`
}

type GRPC struct {
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env:"GRPC_PORT" env-required:"true"`
}

var (
	once           sync.Once
	configInstance *Config
)

func MustLoad() *Config {
	once.Do(func() {
		configInstance = &Config{}
		var err error

		err = tryLoadFromEnv(configInstance)
		if err == nil {
			log.Println("config loaded from environment variables")
			return
		}
		log.Printf("failed to load from env: %v, falling back to config file", err)

		configPath := fetchConfigPath()
		if configPath == "" {
			log.Fatal("config path is empty")
		}

		configInstance, err = LoadPath(configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
	})

	return configInstance
}

func tryLoadFromEnv(cfg *Config) error {
	nodeIDStr := os.Getenv("NODE_ID")
	if nodeIDStr != "" {
		if id, err := strconv.Atoi(nodeIDStr); err == nil {
			cfg.Node.ID = id
		} else {
			return fmt.Errorf("invalid NODE_ID: %v", err)
		}
	}

	portStr := os.Getenv("GRPC_PORT")
	log.Printf("Useful info, grpc port: %s, \n", portStr)
	if strings.HasPrefix(portStr, "${NOMAD_PORT_") {
		cfg.GRPC.Port = 8080
	} else if portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.GRPC.Port = port
		} else {
			return fmt.Errorf("invalid GRPC_PORT: %v", err)
		}
	}

	return cleanenv.ReadEnv(cfg)
}

func LoadPath(configPath string) (*Config, error) {
	// Clean the path to remove any oddities
	cleanPath := filepath.Clean(configPath)

	// Verify the file exists and is accessible
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", cleanPath)
	}

	var cfg Config

	// Read config with cleanenv which also supports environment variables
	err := cleanenv.ReadConfig(cleanPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Validate configuration if needed
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var configPath string

	// Set up flag
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	// If not set by flag, check environment variable
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	// Default config path if none specified
	if configPath == "" {
		configPath = "./config/local.yaml"
	}

	return configPath
}

// validateConfig performs basic configuration validation
func validateConfig(cfg *Config) error {
	if cfg.Node.ID < 0 {
		return fmt.Errorf("node ID must be positive")
	}

	if cfg.GRPC.Host == "" {
		return fmt.Errorf("gRPC host cannot be empty")
	}

	if cfg.GRPC.Port <= 0 || cfg.GRPC.Port > 65535 {
		return fmt.Errorf("gRPC port must be between 1 and 65535")
	}

	return nil
}

func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%d", c.GRPC.Host, c.GRPC.Port)
}
