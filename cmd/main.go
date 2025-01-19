package main

import (
	"Rate_Limiter/internal/config"
	"Rate_Limiter/internal/server"
	"log"
)

func main() {

	loadConfig := config.LoadConfig()

	if err := server.SetupWebServer(loadConfig).Listen(":" + loadConfig.Server.Port); err != nil {
		log.Fatal(err)
	}

}
