package migration

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func Init(dbURL string) error {
	m, err := migrate.New("github://melodyogonna/jobman/migrations", dbURL)
	if err != nil {
		// We don't consider no change as an error
		if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
			return nil
		}
		return err
	}

	err = m.Up()
	if err != nil {
		// We don't consider no change as an error
		if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
			return nil
		}
		return err
	}
	return nil
}

func DeInit(dbURL string) error {
	m, err := migrate.New("github://melodyogonna/jobman/migrations", dbURL)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil {
		return err
	}
	return nil
}
