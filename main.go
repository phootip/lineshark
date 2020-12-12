package main

import (
  "fmt"
	"net/http"
  
  "github.com/phootip/lineshark/config"
	"github.com/phootip/lineshark/server"
)

func main() {
  config, _ := config.InitConfig()
  fmt.Println(config)
  server := server.InitServer()
  
	option := &http.Server{
    Addr: config.Address,
  }
	server.Logger.Fatal(server.StartServer(option))
}
