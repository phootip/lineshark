package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuration enable other app usage
type Configuration struct {
	Port	string `env:"Port"`
}

func init() {
	log.Println("Initializing Config...")
	_ = godotenv.Load()
}

// InitConfig load configuration once
func InitConfig() (*Configuration, error){
	cfg := Configuration{
		os.Getenv("PORT"),
	}

	return &cfg, nil
}
