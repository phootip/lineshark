package main

import (
  "fmt"
	"net/http"
  
  "github.com/phootip/lineshark/config"
	"github.com/phootip/lineshark/server"
	"github.com/phootip/lineshark/controller"
	// "github.com/line/line-bot-sdk-go/linebot"
)

func main() {
  config, _ := config.InitConfig()
  fmt.Println(config)
  server := server.InitServer()
  controller.InitLineBot()

  server.POST("/callback",controller.HandlerCallback)
  controller.InitAPI(server.Group("/api"))
	option := &http.Server{
    Addr: config.Address,
  }
	server.Logger.Fatal(server.StartServer(option))
}
