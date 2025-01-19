package main

import (
	"github.com/valms/fc_rate_limiter/internal/config"
	"github.com/valms/fc_rate_limiter/internal/server"
	"log"
)

func main() {
	loadConfig := config.LoadConfig()

	if err := server.SetupWebServer(loadConfig).Listen(":" + loadConfig.Server.Port); err != nil {
		log.Fatal(err)
	}

}
