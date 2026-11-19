package data

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const (
	dbName     = "bidness"
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbPort     = 5432
)

type DB struct {
	*sql.DB
	pg *embeddedpostgres.EmbeddedPostgres
}

func (db *DB) Close() error {
	return errors.Join(db.DB.Close(), db.pg.Close())
}

func SetupDB(
	ctx context.Context,
	dataFolder string,
	shouldClear bool,
) (*DB, error) {
	dbFolder := filepath.Join(dataFolder, "database")
	dbDataDir := filepath.Join(dbFolder, "postgres")

	// Prepare database directory
	if err := os.MkdirAll(dbFolder, 0755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	// Delete existing DB if shouldClear is true
	if shouldClear {
		if err := os.RemoveAll(dbDataDir); err != nil {
			return nil, fmt.Errorf("remove database directory %s: %w", dbDataDir, err)
		}
	}

	// Start embedded Postgres, creating the data directory if it doesn't exist
	pg := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Database(dbName).
			Username(dbUser).
			Password(dbPassword).
			Port(dbPort).
			Storage(embeddedpostgres.Storage{DataDir: dbDataDir}),
	)
	if err := pg.StartContext(ctx); err != nil {
		return nil, fmt.Errorf("start Postgres: %w", err)
	}

	sqlDB, err := sql.Open("pgx", pg.ConnectionURL())
	if err != nil {
		pg.Close()
		return nil, fmt.Errorf("open Postgres: %w", err)
	}
	db := &DB{DB: sqlDB, pg: pg}

	// Close connection and stop Postgres if setup fails
	ok := false
	defer func() {
		if !ok {
			db.Close()
		}
	}()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("connect to Postgres: %w", err)
	}

	// Apply migrations
	if err := runMigrations(ctx, db.DB); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	// Seed initial data
	if err := SeedDBIfEmpty(
		ctx,
		db.DB,
		filepath.Join(dataFolder, "seed"),
	); err != nil {
		return nil, fmt.Errorf("seed database: %w", err)
	}

	ok = true
	return db, nil
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	files, err := fs.Glob(migrationsFS, "migrations/*sql")
	if err != nil {
		return err
	}

	slices.Sort(files)

	for _, filename := range files {
		content, err := migrationsFS.ReadFile((filename))
		if err != nil {
			return fmt.Errorf("read %s: %w", filename, err)
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("apply %s: %w", filename, err)
		}
	}

	return nil
}
