package sqlite

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Magarate() {
	// Define the database URL for SQLite3
	dbURL := "sqlite://./pkg/db/sqlite/sqlite.db"

	// Define the migration source
	sourceURL := "file://./pkg/db/migrations" // Update with the path to your migrations

	// Initialize a new instance of the Migrate struct
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Apply the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migrations applied successfully.")
}
