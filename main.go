package main

import (
	"log"
	"net/http"

	"github.com/phootip/lineshark/config"
	"github.com/phootip/lineshark/controller"
	"github.com/phootip/lineshark/server"
	// "github.com/line/line-bot-sdk-go/linebot"
)

func main() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  config, _ := config.InitConfig()
  server := server.Server

  server.POST("/callback",controller.HandlerCallback)
  controller.InitAPI(server.Group("/api"))
  controller.InitSheetRoute(server.Group("/sheet"))
	option := &http.Server{
    Addr: config.Address,
  }

  // controller.Temp()

	server.Logger.Fatal(server.StartServer(option))
}
