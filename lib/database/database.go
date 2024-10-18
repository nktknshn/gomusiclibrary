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

func FileFromStatement(stmt *sqlite.Stmt) (models.File, error) {
	var res models.File
	var err error

	res.ID = models.FileID(stmt.GetInt64("id"))
	res.Path = stmt.GetText("path")
	res.Size = stmt.GetInt64("size")
	res.Sha256Hash = stmt.GetText("hash")

	res.Mtime, err = time.Parse(time.DateTime, stmt.GetText("mtime"))

	if err != nil {
		return res, fmt.Errorf("map file: %w", err)
	}

	res.CreatedAt, err = time.Parse(time.DateTime, stmt.GetText("created_at"))

	if err != nil {
		return res, fmt.Errorf("map file: %w", err)
	}

	idxUpdateAt := stmt.ColumnIndex("updated_at")

	if stmt.ColumnType(idxUpdateAt) == sqlite.TypeNull {
		res.UpdatedAt = time.Time{}
	} else {
		res.UpdatedAt, err = time.Parse(time.DateTime, stmt.ColumnText(idxUpdateAt))
		if err != nil {
			return res, fmt.Errorf("map file: %w", err)
		}
	}

	idxDeletedAt := stmt.ColumnIndex("deleted_at")

	if stmt.ColumnType(idxDeletedAt) == sqlite.TypeNull {
		res.DeletedAt = time.Time{}
	} else {
		res.DeletedAt, err = time.Parse(time.DateTime, stmt.ColumnText(idxDeletedAt))
		if err != nil {
			return res, fmt.Errorf("map file: %w", err)
		}
	}

	return res, nil
}

func (db *Database) FilesList() (models.FileSlice, error) {

	var result []models.File

	err := sqlitex.Execute(db.conn, "SELECT id, path, size, hash, mtime, created_at, updated_at, deleted_at FROM file;", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			f, err := FileFromStatement(stmt)

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
INSERT INTO file (path, size, hash, mtime, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
*/
func (db *Database) FilesInsert(files models.FileSlice) error {

	for _, f := range files {

		stmt, err := db.conn.Prepare("INSERT INTO file (path, size, hash, mtime, created_at, updated_at, deleted_at) VALUES ($path, $size, $hash, $mtime, $created_at, $updated_at, $deleted_at);")

		if err != nil {
			return fmt.Errorf("database: files insert: %w", err)
		}

		stmt.SetText("$path", f.Path)
		stmt.SetInt64("$size", f.Size)
		stmt.SetText("$hash", f.Sha256Hash)
		stmt.SetText("$mtime", f.Mtime.Format(time.DateTime))

		if f.CreatedAt.IsZero() {
			stmt.SetText("$created_at", time.Now().Format(time.DateTime))
		} else {
			stmt.SetText("$created_at", f.CreatedAt.Format(time.DateTime))
		}

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
	}

	return nil
}
