package server

import (
	"github.com/phootip/lineshark/config"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
  
)

var (
	// Server sigular
	Server *echo.Echo
)

// InitServer initialize once
func InitServer() *echo.Echo{
	fmt.Println("Initialize server")
	
	Server = echo.New()
	Server.Use(middleware.Logger())
	Server.Use(middleware.Recover())
	
	Server.GET("/", hello)
	api := Server.Group("/api")
  api.GET("/phootip", phootip)
	return Server
}

// Start to start server
func Start(config *config.Configuration) {
	option := &http.Server{
    Addr: config.Address,
  }
	Server.Logger.Fatal(Server.StartServer(option))
}

// Handler
func hello(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, World!")
}

// Handler
func phootip(c echo.Context) error {
  return c.String(http.StatusOK, "Hello, phootip!")
}
