package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuration enable other app usage
type Configuration struct {
	Address	string `env:"ADDRESS"`
}

func init() {
	log.Println("Initializing Config...")
	_ = godotenv.Load()
}

// InitConfig load configuration once
func InitConfig() (*Configuration, error){
	cfg := Configuration{
		os.Getenv("ADDRESS"),
	}

	return &cfg, nil
}
