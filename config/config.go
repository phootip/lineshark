package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Address	string `env:"ADDRESS"`
}

func InitConfig() (*Configuration, error){
	fmt.Println("Initializing Config...")
	_ := godotenv.Load()
}
