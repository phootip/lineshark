package config

import (
	"fmt"
	"os"

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

	cfg := Configuration{
		os.Getenv("ADDRESS"),
	}

	return &cfg, nil
}
