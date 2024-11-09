package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// const (
// 	// DB_HOST     = "postgres"  // for docker swarm
// 	// DB_HOST = "host.minikube.internal" // for kubernetes
// 	DB_HOST = "192.168.49.2" // for kubernetes

// 	DB_PORT     = "5432"
// 	DB_USER     = "postgres"
// 	DB_PASSWORD = "root"
// 	DB_NAME     = "go_microservice_auth_service"
// )

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func LoadEnv() (*Config, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	log.Print("host", host)
	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		return nil, fmt.Errorf("missing config values")
	}

	config := &Config{
		DBHost: host,
		DBPort: port,
		DBUser: user,
		DBPass: password,
		DBName: dbName,
	}

	// config := &Config{
	// 	DBHost: DB_HOST,
	// 	DBPort: DB_PORT,
	// 	DBUser: DB_USER,
	// 	DBPass: DB_PASSWORD,
	// 	DBName: DB_NAME,
	// }
	return config, nil
}
