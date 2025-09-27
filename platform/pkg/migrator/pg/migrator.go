package pg

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func New(db *sql.DB, migrationsPath string) (*Migrator, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &Migrator{
		migrate: m,
	}, nil
}

func (m *Migrator) Up() error {
	return m.migrate.Up()
}

func (m *Migrator) Down() error {
	return m.migrate.Down()
}

func (m *Migrator) Force(version int) error {
	return m.migrate.Force(version)
}

func (m *Migrator) Version() (int, bool, error) {
	version, dirty, err := m.migrate.Version()
	return int(version), dirty, err
}

func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}
