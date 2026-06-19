package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Zrossiz/finance-backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := log.Default()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("migrator get config err: ", err)
	}

	var dbHost string
	flag.StringVar(&dbHost, "h", "localhost", "address and port to run server")
	flag.Parse()

	uri := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		dbHost,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	)

	m, err := migrate.New("file://migrations", uri)
	if err != nil {
		logger.Fatal("init migrate err: ", err)
	}
	m.Steps(2)

	logger.Print("migration successful\n")
}
