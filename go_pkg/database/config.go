package database

import (
	"fmt"
	"log"
	"os"
)

// Config holds database connection parameters
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConfig creates a database config from environment variables
func NewConfig() Config {
	return Config{
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     getEnvWithDefault("DB_PORT", "3306"),
		User:     getEnvWithDefault("DB_USER", "maria"),
		Password: getEnvWithDefault("DB_PW", "maria123"),
		DBName:   getEnvWithDefault("DB_NAME", "mariadb"),
	}
}

// DSN returns a MySQL connection string
func (c *Config) DSN() string {
	log.Println(c.User, c.Password, c.Host, c.Port, c.DBName)
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
