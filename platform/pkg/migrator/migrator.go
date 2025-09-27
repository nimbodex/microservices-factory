package migrator

import (
	"context"
	"database/sql"
)

type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
	Force(ctx context.Context, version int) error
	Version(ctx context.Context) (int, bool, error)
}

type migrator struct {
	db      *sql.DB
	source  Source
	table   string
	version int
	dirty   bool
}

type Source interface {
	Close() error
	Open(url string) (Source, error)
	Up() error
	Down() error
	Force(version int) error
	Version() (int, bool, error)
}

func New(db *sql.DB, source Source, table string) Migrator {
	return &migrator{
		db:     db,
		source: source,
		table:  table,
	}
}

func (m *migrator) Up(ctx context.Context) error {
	return m.source.Up()
}

func (m *migrator) Down(ctx context.Context) error {
	return m.source.Down()
}

func (m *migrator) Force(ctx context.Context, version int) error {
	return m.source.Force(version)
}

func (m *migrator) Version(ctx context.Context) (int, bool, error) {
	return m.source.Version()
}
