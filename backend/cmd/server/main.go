package main

import (
	"log"
	"os"

	"antd-gin-admin-backend/internal/server"
	"antd-gin-admin-backend/internal/config"
)

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "configs/config.yaml"
	}
	
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	if err := server.Run(cfg); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

