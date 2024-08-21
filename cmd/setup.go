package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MustMigrate(dbAddr, migrationUrl string) {
	m, err := migrate.New(
		migrationUrl,
		dbAddr,
	)
	if err != nil {
		log.Fatal("Failed to create a migration: ", err)
	}
	defer m.Close()
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not migrate up: %s", err.Error())
	}
	log.Println("Migrated successfully")
}
