package main

import (
	"fmt"
	"log"
	"os"

	"key-distribution-system/internal/config"
	"key-distribution-system/internal/db"
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

	sqlitePath := cfg.SQLite.Path
	if sqlitePath == "" {
		sqlitePath = "data/kds.db"
	}

	if err = db.Init(sqlitePath); err != nil {
		log.Fatalf("init db failed: %v", err)
	}

	store.Init(db.DB)

	r := router.New()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err = r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
