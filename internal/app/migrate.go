package app

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (a *App) Migrate() {
	fmt.Println("\n\n\n### ### ### MIGRATE ### ### ###\n\n\n")
	driver, err := postgres.WithInstance(a.DB.DB, &postgres.Config{})
	if err != nil {
		log.Println(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations/", "pmdbstore", driver)
	if err != nil {
		log.Printf("[App.Migrate]: %e", err)
	}
	if err := m.Steps(2); err != nil {
		log.Println(err)
	}
	fmt.Println("\n\n\n### ### ### MIGRATION FINISHED ### ### ###\n\n\n")
}
