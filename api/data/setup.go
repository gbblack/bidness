package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	_ "github.com/duckdb/duckdb-go/v2"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func SetupDB(
	ctx context.Context,
	dataFolder string,
	shouldClear bool,
) (*sql.DB, error) {
	dbFolder := filepath.Join(dataFolder, "database")
	dbFilename := filepath.Join(dbFolder, "bidness.duckdb")

	// Prepare database directory
	if err := os.MkdirAll(dbFolder, 0755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	// Delete existing DB is shouldClear is true
	if shouldClear {
		for _, filename := range []string{
			dbFilename,
			dbFilename + ".wal",
		} {
			if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
				return nil, fmt.Errorf("remove database file %s: %w", filename, err)
			}
		}
	}

	// Create and open DuckDB database file
	db, err := sql.Open("duckdb", dbFilename)
	if err != nil {
		return nil, fmt.Errorf("open DuckDB: %w", err)
	}

	// Close connection if setup fails
	ok := false
	defer func() {
		if !ok {
			db.Close()
		}
	}()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("connect to DuckDB: %w", err)
	}

	// Apply migrations
	if err := runMigrations(ctx, db); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	// Seed initial data
	if err := SeedDBIfEmpty(
		ctx,
		db,
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
