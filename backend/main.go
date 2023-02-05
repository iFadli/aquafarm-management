package main

import (
	"aquafarm-management/app"
	"aquafarm-management/app/config"
	_ "aquafarm-management/docs"
	"log"
)

// @title           AquaFarm Management API
// @version         0.6
// @description     This is Assignment Workspace for Coding Test DELOS Aqua.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Fadli Zul Fahmi
// @contact.url    http://www.linkedin.com/in/fadli-zul-fahmi
// @contact.email  hi@fadli.dev

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apiKey APIKeyHeader
// @in header
// @name Authorization

func main() {
	// Load config
	cfg := config.Load()

	// Initialize app
	apps := app.New(cfg)

	// Run server
	if err := apps.Run(cfg); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
