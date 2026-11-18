package data

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
)

func SeedDBIfEmpty(
	ctx context.Context,
	db *sql.DB,
	dataFolder string,
) error {
	var count int64

	if err := db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM authors",
	).Scan(&count); err != nil {
		return fmt.Errorf("failed to count authors: %w", err)
	}

	if count > 0 {
		return nil
	}

	filename := filepath.Join(dataFolder, "authors.json")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin seed transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO authors (
			key,
			name, type_key,
			source_records,
			latest_revision,
			revision,
			created,
			last_modified
		)
		SELECT
			key,
			name,
			type.key,
			source_records,
			latest_revision,
			revision,
			CAST(created.value AS TIMESTAMP),
			CAST(last_modified.value AS TIMESTAMP)
		FROM read_json_auto(?)
	`, filename)
	if err != nil {
		return fmt.Errorf("failed to import authors: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit authros seed: %w", err)
	}

	return nil
}
