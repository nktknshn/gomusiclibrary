package database

import (
	"fmt"
	"time"

	"github.com/nktknshn/gomusiclibrary/lib/database/models"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Database struct {
	conn *sqlite.Conn
}

func New(conn *sqlite.Conn) *Database {
	return &Database{conn}
}

func NewFromFile(path string) (*Database, error) {
	conn, err := sqlite.OpenConn(path, sqlite.OpenReadWrite)

	if err != nil {
		return nil, fmt.Errorf("database: create connection: %w", err)
	}

	return New(conn), nil
}

func (db *Database) FilesList() (models.FileSlice, error) {

	var result []models.File

	err := sqlitex.Execute(db.conn, "SELECT id, path, size, hash, ctime, mtime, created_at, updated_at, deleted_at FROM file;", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			// fmt.Println(stmt)
			var err error
			var f models.File
			f.ID = models.FileID(stmt.ColumnInt64(0))
			f.Path = stmt.ColumnText(1)
			f.Size = stmt.ColumnInt64(2)
			f.Sha256Hash = stmt.ColumnText(3)
			f.Ctime = stmt.ColumnInt64(4)
			f.Mtime = stmt.ColumnInt64(5)

			f.CreatedAt, err = time.Parse(time.DateTime, stmt.ColumnText(6))

			if err != nil {
				return fmt.Errorf("database: list files: %w", err)
			}

			if stmt.ColumnType(7) == sqlite.TypeNull {
				f.UpdatedAt = time.Time{}
			} else {
				f.UpdatedAt, err = time.Parse(time.DateTime, stmt.ColumnText(7))
				if err != nil {
					return fmt.Errorf("database: list files: %w", err)
				}
			}

			if stmt.ColumnType(8) == sqlite.TypeNull {
				f.DeletedAt = time.Time{}
			} else {
				f.DeletedAt, err = time.Parse(time.DateTime, stmt.ColumnText(8))
				if err != nil {
					return fmt.Errorf("database: list files: %w", err)
				}
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
INSERT INTO file (path, size, hash, ctime, mtime, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
*/
func (db *Database) FilesInsert(files models.FileSlice) error {

	for _, f := range files {

		stmt, err := db.conn.Prepare("INSERT INTO file (path, size, hash, ctime, mtime, created_at, updated_at, deleted_at) VALUES ($path, $size, $hash, $ctime, $mtime, $created_at, $updated_at, $deleted_at);")

		if err != nil {
			return fmt.Errorf("database: files insert: %w", err)
		}

		stmt.SetText("$path", f.Path)
		stmt.SetInt64("$size", f.Size)
		stmt.SetText("$hash", f.Sha256Hash)
		stmt.SetInt64("$ctime", f.Ctime)
		stmt.SetInt64("$mtime", f.Mtime)
		stmt.SetText("$created_at", f.CreatedAt.Format(time.DateTime))

		if !f.UpdatedAt.IsZero() {
			stmt.SetText("$updated_at", f.UpdatedAt.Format(time.DateTime))
		}
		if !f.DeletedAt.IsZero() {
			stmt.SetText("$deleted_at", f.DeletedAt.Format(time.DateTime))
		}

		_, err = stmt.Step()

		if err != nil {
			return fmt.Errorf("database: files insert: %w", err)
		}

		// Prepare clears it for us
		// if err := stmt.Reset(); err != nil {
		// 	return fmt.Errorf("database: files insert: %w", err)
		// }

	}

	return nil
}
