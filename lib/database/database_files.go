package database

import (
	"fmt"

	"github.com/nktknshn/gomusiclibrary/lib/database/models"
	"github.com/nktknshn/gomusiclibrary/lib/util/timeutil"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func (db *Database) FilesList() (models.FileSlice, error) {

	var result []models.File

	query := "SELECT id, path, size, hash, mtime, created_at, updated_at, deleted_at FROM file;"

	err := sqlitex.Execute(db.conn, query, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			f, err := fileFromStatement(stmt)
			if err != nil {
				return err
			}
			result = append(result, f)
			return nil
		},
	})

	if err != nil {
		return nil, fmt.Errorf("database: list files: %w", err)
	}

	return result, nil
}

/*
Refreshing library.

- If a file exists
	- if it needs to be updated: mtime or size changed
	- If nothing changed, update last check date (?)

- If a file does not exist, insert it

- We must check if a file is deleted
	- Mountpoint might be missing
	- File might be deleted

- When a file updated, we must update tags. Which means:
	1. Remove removed tags
	2. Update updated values
*/

// Fails if a path already exists
func (db *Database) FilesInsert(files models.FileSlice) error {

	t := timeutil.Measurer()

	// _ = db.conn.Prep("PRAGMA synchronous = OFF;")

	query := "INSERT INTO file (path, size, hash, mtime, created_at, updated_at, deleted_at) VALUES ($path, $size, $hash, $mtime, $created_at, $updated_at, $deleted_at);"

	for _, f := range files {

		stmt, err := db.conn.Prepare(query)

		if err != nil {
			return fmt.Errorf("database: files insert: %w", err)
		}

		SetFileToStatement(stmt, &f)

		_, err = stmt.Step()

		if err != nil {
			return fmt.Errorf("database: files insert: %w", err)
		}
	}

	t.Print("Inserting")
	// fmt.Printf("Inserting %d took %.2f secs\n", len(files), time.Since(timeNow).Seconds())

	return nil
}
