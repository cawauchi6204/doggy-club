package main

import (
	"flag"
	"log"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/db"
)

func main() {
	var seed bool
	flag.BoolVar(&seed, "seed", false, "Seed initial data")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	database, err := db.InitPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.AutoMigrate(database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Seed data if requested
	if seed {
		if err := db.SeedInitialData(database); err != nil {
			log.Fatal("Failed to seed data:", err)
		}
	}

	log.Println("Migration completed successfully")
}