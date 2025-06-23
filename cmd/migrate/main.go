package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func main() {
	// Parse command line arguments
	var direction string
	flag.StringVar(&direction, "direction", "up", "Migration direction (up or down)")
	flag.Parse()

	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Get database configuration
	dbConfig := viper.GetStringMapString("database.postgres")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig["user"],
		dbConfig["password"],
		dbConfig["host"],
		dbConfig["port"],
		dbConfig["dbname"],
		dbConfig["sslmode"],
	)

	// Initialize migrator
	migrationsPath := filepath.Join("migrations")
	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		log.Fatalf("Error creating migrator: %s", err)
	}
	defer m.Close()

	// Run migrations
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error running migrations: %s", err)
		}
		log.Println("Successfully ran migrations")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error reverting migrations: %s", err)
		}
		log.Println("Successfully reverted migrations")
	default:
		log.Fatalf("Invalid direction: %s", direction)
	}
}
