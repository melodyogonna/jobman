package migration

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func Init(dbURL string) error {
	m, err := migrate.New("github://melodyogonna/jobman/migrations", dbURL)
	if err != nil {
		// We don't consider no change as an error
		if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
			log.Print(err)
			return nil
		}
		return err
	}

	return m.Up()
}

func DeInit(dbURL string) error {
	m, err := migrate.New("github://melodyogonna/jobman/migrations", dbURL)
	if err != nil {
		// We don't consider no change as an error
		if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
			log.Print(err)
			return nil
		}
		return err
	}

	return m.Down()
}
