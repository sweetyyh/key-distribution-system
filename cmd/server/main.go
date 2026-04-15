package main

import (
	"fmt"
	"log"
	"os"

	"key-distribution-system/internal/config"
	"key-distribution-system/internal/router"
	"key-distribution-system/internal/store"
)

func main() {
	configPath := os.Getenv("KDS_CONFIG")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	store.Init()

	r := router.New()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err = r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
