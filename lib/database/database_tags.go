package database

import (
	"fmt"
	"time"

	"github.com/nktknshn/gomusiclibrary/lib/database/models"
	"github.com/nktknshn/gomusiclibrary/lib/util/sqlutil"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func (db *Database) TagsList() (models.AudioTagSlice, error) {
	var result []models.AudioTag

	query := "SELECT id, file_id, name, value, created_at, updated_at, deleted_at FROM file_audio_tag;"

	err := sqlitex.Execute(db.conn, query, &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			t, err := audioTagFromStatement(stmt)
			if err != nil {
				return err
			}
			result = append(result, t)
			return nil
		},
	})

	if err != nil {
		return nil, fmt.Errorf("database: list tags: %w", err)
	}

	return result, nil
}

func (db *Database) TagsInsert(tags models.AudioTagSlice) error {
	query := "INSERT INTO file_audio_tag (file_id, name, value, created_at, updated_at, deleted_at) VALUES ($file_id, $name, $value, $created_at, $updated_at, $deleted_at);"

	for _, t := range tags {
		stmt, err := db.conn.Prepare(query)

		if err != nil {
			return fmt.Errorf("database: tags insert: %w", err)
		}

		setAudioTagToStatement(stmt, &t)

		_, err = stmt.Step()

		if err != nil {
			return fmt.Errorf("database: tags insert: %w", err)
		}

	}

	// reset stm?
	return nil
}

func setAudioTagToStatement(stmt *sqlite.Stmt, t *models.AudioTag) {
	stmt.SetInt64("$file_id", int64(t.FileID))
	stmt.SetText("$name", t.Name)
	stmt.SetText("$value", t.Value)

	sqlutil.SetDateTextOr(stmt, "$created_at", t.CreatedAt, time.Now())
	sqlutil.SetDateTextIfNotZero(stmt, "$updated_at", t.UpdatedAt)
	sqlutil.SetDateTextIfNotZero(stmt, "$deleted_at", t.DeletedAt)
}

func audioTagFromStatement(stmt *sqlite.Stmt) (models.AudioTag, error) {
	var res models.AudioTag
	var err error

	res.FileID = models.FileID(stmt.GetInt64("file_id"))
	res.Name = stmt.GetText("name")
	res.Value = stmt.GetText("value")

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
