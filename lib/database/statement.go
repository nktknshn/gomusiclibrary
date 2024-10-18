package database

import (
	"fmt"
	"time"

	"github.com/nktknshn/gomusiclibrary/lib/database/models"
	"github.com/nktknshn/gomusiclibrary/lib/util/sqlutil"
	"zombiezen.com/go/sqlite"
)

func SetFileToStatement(stmt *sqlite.Stmt, f *models.File) {
	stmt.SetText("$path", f.Path)
	stmt.SetInt64("$size", f.Size)
	stmt.SetText("$hash", f.Sha256Hash)
	stmt.SetText("$mtime", f.Mtime.Format(time.DateTime))
	sqlutil.SetDateTextOr(stmt, "$created_at", f.CreatedAt, time.Now())
	sqlutil.SetDateTextIfNotZero(stmt, "$updated_at", f.UpdatedAt)
	sqlutil.SetDateTextIfNotZero(stmt, "$deleted_at", f.DeletedAt)

	// if f.CreatedAt.IsZero() {
	// 	stmt.SetText("$created_at", time.Now().Format(time.DateTime))
	// } else {
	// 	stmt.SetText("$created_at", f.CreatedAt.Format(time.DateTime))
	// }

	// if !f.UpdatedAt.IsZero() {
	// 	stmt.SetText("$updated_at", f.UpdatedAt.Format(time.DateTime))
	// }
	// if !f.DeletedAt.IsZero() {
	// 	stmt.SetText("$deleted_at", f.DeletedAt.Format(time.DateTime))
	// }
}

func fileFromStatement(stmt *sqlite.Stmt) (models.File, error) {
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

	res.UpdatedAt, err = sqlutil.DateFromColumnMaybeNull(stmt, "updated_at")

	if err != nil {
		return res, fmt.Errorf("map file: %w", err)
	}

	res.DeletedAt, err = time.Parse(time.DateTime, "deleted_at")

	if err != nil {
		return res, fmt.Errorf("map file: %w", err)
	}

	return res, nil
}
