package main

import (
	"log"

	"github.com/gadzz/scopepay/internal/config"
	"github.com/gadzz/scopepay/internal/db"
	"github.com/gadzz/scopepay/internal/router"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	r := router.New(cfg, database)

	log.Printf("listening on :%s (db=%s)", cfg.Port, cfg.DBDriver)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
