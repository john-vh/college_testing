package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/john-vh/college_testing/backend/env"
)

var EXT_ENVIRONMENT string = env.DEV

func main() {
	cfg := env.GetConfig(EXT_ENVIRONMENT)

	databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.POSTGRES_USER, cfg.POSTGRES_PASSWORD, cfg.POSTGRES_HOST, cfg.POSTGRES_PORT, cfg.POSTGRES_DB)
	migrationsPath := "file:///migrations"

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "force") {
		log.Fatal("Expects either 'up', 'down', or 'drop' as second and only argument")
	}

	cmd := os.Args[1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		return
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		return
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatal(err)
		}
		return
	case "force":
		if val, err := strconv.Atoi(os.Args[2]); err == nil {
			if err := m.Force(val); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	if val, err := strconv.Atoi(cmd); err == nil {
		if err := m.Steps(val); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.Fatal("Expected 'up', 'down', or 'drop', got: " + cmd)
}
