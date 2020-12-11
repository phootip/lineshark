package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

// Configuration enable other app usage
type Configuration struct {
	Address	string `env:"ADDRESS"`
}

// InitConfig load configuration once
func InitConfig() (*Configuration, error){
	fmt.Println("Initializing Config...")
	_ = godotenv.Load()

	cfg := Configuration{}

	return &cfg, nil
}
