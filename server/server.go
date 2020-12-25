package server

import (
	"github.com/phootip/lineshark/config"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
  
)

var (
	// Server sigular
	Server *echo.Echo
)

func init() {
	log.Println("Initialize server")
	
	Server = echo.New()
	Server.Use(middleware.Logger())
	Server.Use(middleware.Recover())
	
	Server.GET("/", hello)
}

// Start to start server
func Start(config *config.Configuration) {
	option := &http.Server{
    Addr: config.Port,
  }
	Server.Logger.Fatal(Server.StartServer(option))
}

// Handler
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}
