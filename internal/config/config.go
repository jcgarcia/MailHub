package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port     string
	LogLevel string

	// Database (SQLite)
	DatabasePath string

	// SSH Configuration
	SSH SSHConfig

	// Auth
	DevMode      bool
	DevAuthEmail string
}

// SSHConfig holds SSH connection settings
type SSHConfig struct {
	Host        string
	Port        int
	User        string
	KeyPath     string
	JumpHost    string
	JumpUser    string
	JumpKeyPath string
}

// Load reads configuration from environment variables
func Load() *Config {
	sshPort, _ := strconv.Atoi(getEnv("CMH_SSH_PORT", "22"))

	return &Config{
		Port:         getEnv("PORT", "8080"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		DatabasePath: getEnv("DATABASE_PATH", "/data/mailhub.db"),

		SSH: SSHConfig{
			Host:        getEnv("CMH_SSH_HOST", "localhost"),
			Port:        sshPort,
			User:        getEnv("CMH_SSH_USER", "postman"),
			KeyPath:     getEnv("CMH_SSH_KEY_PATH", "/secrets/mailhub_key"),
			JumpHost:    getEnv("CMH_SSH_JUMP_HOST", "jump.ingasti.com"),
			JumpUser:    getEnv("CMH_SSH_JUMP_USER", "ubuntu"),
			JumpKeyPath: getEnv("CMH_SSH_JUMP_KEY_PATH", "/secrets/jump_key"),
		},

		DevMode:      getEnv("DEV_MODE", "false") == "true",
		DevAuthEmail: getEnv("DEV_AUTH_EMAIL", "dev@example.com"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
