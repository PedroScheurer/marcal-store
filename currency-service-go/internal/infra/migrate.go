package infra

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseUrl string) error {
	if !strings.HasPrefix(databaseUrl, "postgres://") {
		return fmt.Errorf("invalid database url: %s", databaseUrl)
	}

	m, err := migrate.New("file://db_migration", databaseUrl)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && errors.Is(err, migrate.ErrNoChange) {
	}

	return nil
}
