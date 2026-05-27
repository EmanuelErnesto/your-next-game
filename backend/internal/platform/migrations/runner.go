package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ApplyUp(ctx context.Context, db *sql.DB, dir string) error {
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".up.sql") {
			files = append(files, filepath.Join(dir, name))
		}
	}
	sort.Strings(files)

	for _, filePath := range files {
		if err := applySingleFile(ctx, db, filePath); err != nil {
			return err
		}
	}

	return nil
}

func applySingleFile(ctx context.Context, db *sql.DB, filePath string) error {
	fileName := filepath.Base(filePath)
	applied, err := isApplied(ctx, db, fileName)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read migration file %q: %w", filePath, err)
	}

	statement := strings.TrimSpace(string(sqlBytes))
	if statement == "" {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx for %q: %w", filePath, err)
	}

	if _, err := tx.ExecContext(ctx, statement); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("execute migration %q: %w", filePath, err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`insert into schema_migrations (filename, applied_at) values ($1, now())`,
		fileName,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %q: %w", filePath, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %q: %w", filePath, err)
	}

	return nil
}

func EnsureDirExists(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("stat migrations dir: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", dir)
	}
	return nil
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
create table if not exists schema_migrations (
  filename text primary key,
  applied_at timestamptz not null
)`)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations table: %w", err)
	}
	return nil
}

func isApplied(ctx context.Context, db *sql.DB, fileName string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(
		ctx,
		`select exists (select 1 from schema_migrations where filename = $1)`,
		fileName,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check migration status for %q: %w", fileName, err)
	}
	return exists, nil
}

func CountUpFiles(dir string) (int, error) {
	count := 0
	err := filepath.WalkDir(dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".up.sql") {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}
